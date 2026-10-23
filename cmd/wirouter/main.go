package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/MathiasPiazzolla/WiRouter/internal/autostart"
	"github.com/MathiasPiazzolla/WiRouter/internal/hotspot"
	"github.com/MathiasPiazzolla/WiRouter/internal/instance"
	"github.com/MathiasPiazzolla/WiRouter/internal/notify"
	"github.com/MathiasPiazzolla/WiRouter/internal/runner"
	"github.com/MathiasPiazzolla/WiRouter/model"
)

const (
	instanceName     = `Local\WiRouter`
	autostartExeName = "wirouter-autostart.exe"
	// All'avvio la connessione può comparire con ritardo, per questo si attende a lungo.
	sourceTimeout   = 5 * time.Minute
	sourcePoll      = 5 * time.Second
	startAttempts   = 3
	startRetryDelay = 10 * time.Second
)

func main() {
	flags := flag.NewFlagSet("wirouter", flag.ContinueOnError)
	autostartMode := flags.Bool("autostart", false, "avvia l'hotspot con i dati dell'attività pianificata, senza interazione")
	installMode := flags.Bool("install-autostart", false, "crea l'attività pianificata che avvia WiRouter all'accesso")
	uninstallMode := flags.Bool("uninstall-autostart", false, "rimuove l'attività pianificata")
	config := model.HotspotConfig{}
	flags.StringVar(&config.Source, "source", "", "connessione da condividere (con --autostart)")
	flags.StringVar(&config.SSID, "ssid", "", "SSID dell'hotspot (con --autostart)")
	flags.StringVar(&config.Password, "password", "", "password dell'hotspot (con --autostart)")
	if err := flags.Parse(os.Args[1:]); err != nil {
		os.Exit(2)
	}

	var err error
	switch {
	case *autostartMode:
		// L'eseguibile senza console non ha dove mostrare errori: usa una finestra.
		if err := runAutostart(config); err != nil {
			_ = notify.Error("WiRouter", err.Error())
			os.Exit(1)
		}
		return
	case *installMode:
		err = installAutostart(os.Stdin, os.Stdout)
	case *uninstallMode:
		err = uninstallAutostart(os.Stdout)
	default:
		err = run(os.Stdin, os.Stdout)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Errore:", err)
		os.Exit(1)
	}
}

func run(input io.Reader, output io.Writer) error {
	release, err := instance.Acquire(instanceName)
	if err != nil {
		return err
	}
	defer release()

	sources, err := availableSources()
	if err != nil {
		return err
	}

	fmt.Fprintln(output, "========================================")
	fmt.Fprintln(output, "       WiRouter - Hotspot Manager")
	fmt.Fprintln(output, "========================================")
	fmt.Fprintln(output)
	printSources(output, sources)

	reader := bufio.NewReader(input)
	config, err := askConfig(reader, output, sources)
	if err != nil {
		return err
	}
	ssid := config.SSID
	source := config.Source

	// Registrato prima dell'avvio: CTRL+C non deve terminare il processo senza fermare l'hotspot.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	if err := hotspot.StartHotspot(config); err != nil {
		return err
	}

	fmt.Fprintln(output)
	fmt.Fprintf(output, "Hotspot %q avviato, condivide la connessione %q.\n", ssid, source)
	fmt.Fprintln(output, "Premi CTRL+C per fermare l'hotspot.")
	<-interrupt

	fmt.Fprintln(output)
	fmt.Fprintln(output, "Arresto hotspot...")
	if err := hotspot.StopHotspot(source); err != nil {
		return err
	}
	fmt.Fprintln(output, "Hotspot fermato.")
	return nil
}

func availableSources() ([]string, error) {
	sources, err := hotspot.ListSources()
	if err != nil {
		return nil, err
	}
	if len(sources) == 0 {
		return nil, errors.New("nessuna connessione con accesso a Internet da condividere")
	}
	return sources, nil
}

func printSources(output io.Writer, sources []string) {
	fmt.Fprintln(output, "Connessioni con accesso a Internet:")
	for index, source := range sources {
		fmt.Fprintf(output, "  %d) %s\n", index+1, source)
	}
}

func askConfig(reader *bufio.Reader, output io.Writer, sources []string) (model.HotspotConfig, error) {
	source, err := readSource(reader, output, sources)
	if err != nil {
		return model.HotspotConfig{}, err
	}
	ssid, err := readValidated(reader, output, "Inserisci SSID hotspot: ", hotspot.ValidateSSID)
	if err != nil {
		return model.HotspotConfig{}, err
	}
	password, err := readValidated(
		reader,
		output,
		"Inserisci password (8-63 caratteri): ",
		hotspot.ValidatePassword,
	)
	if err != nil {
		return model.HotspotConfig{}, err
	}
	return model.HotspotConfig{SSID: ssid, Password: password, Source: source}, nil
}

// runAutostart è la modalità usata dall'attività pianificata. Si ferma con
// CTRL+C o con il segnale di logout/arresto, e in quel caso spegne l'hotspot.
func runAutostart(config model.HotspotConfig) error {
	release, err := instance.Acquire(instanceName)
	if err != nil {
		return err
	}
	defer release()

	if config.Source == "" {
		return errors.New("manca la connessione da condividere (--source)")
	}
	if err := hotspot.ValidateSSID(config.SSID); err != nil {
		return err
	}
	if err := hotspot.ValidatePassword(config.Password); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	err = runner.Run(ctx, hotspot.Controller{}, runner.Options{
		Config:          config,
		SourceTimeout:   sourceTimeout,
		SourcePoll:      sourcePoll,
		StartAttempts:   startAttempts,
		StartRetryDelay: startRetryDelay,
	})
	if errors.Is(err, context.Canceled) {
		return nil
	}
	return err
}

func installAutostart(input io.Reader, output io.Writer) error {
	sources, err := availableSources()
	if err != nil {
		return err
	}
	printSources(output, sources)

	reader := bufio.NewReader(input)
	config, err := askConfig(reader, output, sources)
	if err != nil {
		return err
	}

	exe, err := os.Executable()
	if err != nil {
		return err
	}
	autostartExe := filepath.Join(filepath.Dir(exe), autostartExeName)
	if _, err := os.Stat(autostartExe); err != nil {
		return fmt.Errorf("%s non trovato accanto a %s: compila anche la versione senza console", autostartExeName, exe)
	}

	userID := autostart.CurrentUserID()
	err = autostart.Install(autostart.Options{
		ExePath: autostartExe,
		UserID:  userID,
		Args:    autostartArgs(config),
	})
	if err != nil {
		return err
	}
	fmt.Fprintf(output, "Attività %q creata: WiRouter partirà all'accesso di %s.\n", autostart.TaskName, userID)
	return nil
}

func uninstallAutostart(output io.Writer) error {
	if err := autostart.Uninstall(); err != nil {
		return err
	}
	fmt.Fprintln(output, "Attività di avvio automatico rimossa.")
	return nil
}

// autostartArgs costruisce gli argomenti dell'attività pianificata. Ogni valore
// è quotato secondo le regole di Windows, così spazi e virgolette restano intatti.
func autostartArgs(config model.HotspotConfig) string {
	return strings.Join([]string{
		"-autostart",
		"-source", syscall.EscapeArg(config.Source),
		"-ssid", syscall.EscapeArg(config.SSID),
		"-password", syscall.EscapeArg(config.Password),
	}, " ")
}

func readSource(reader *bufio.Reader, output io.Writer, sources []string) (string, error) {
	for {
		value, err := readLine(reader, output, "Seleziona la connessione da condividere (numero): ")
		if err != nil {
			return "", err
		}

		index, err := strconv.Atoi(value)
		if err != nil || index < 1 || index > len(sources) {
			fmt.Fprintf(output, "Scelta non valida: inserisci un numero tra 1 e %d.\n", len(sources))
			continue
		}
		return sources[index-1], nil
	}
}

func readValidated(
	reader *bufio.Reader,
	output io.Writer,
	prompt string,
	validate func(string) error,
) (string, error) {
	for {
		value, err := readLine(reader, output, prompt)
		if err != nil {
			return "", err
		}
		if err := validate(value); err != nil {
			fmt.Fprintln(output, "Valore non valido:", err)
			continue
		}
		return value, nil
	}
}

func readLine(reader *bufio.Reader, output io.Writer, prompt string) (string, error) {
	fmt.Fprint(output, prompt)
	value, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("lettura input: %w", err)
	}

	value = strings.TrimSpace(value)
	if errors.Is(err, io.EOF) && value == "" {
		return "", errors.New("input terminato")
	}
	return value, nil
}
