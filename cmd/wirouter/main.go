package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/MathiasPiazzolla/WiRouter/internal/hotspot"
	"github.com/MathiasPiazzolla/WiRouter/model"
)

func main() {
	if err := run(os.Stdin, os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "Errore:", err)
		os.Exit(1)
	}
}

func run(input io.Reader, output io.Writer) error {
	sources, err := hotspot.ListSources()
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return errors.New("nessuna connessione con accesso a Internet da condividere")
	}

	fmt.Fprintln(output, "========================================")
	fmt.Fprintln(output, "       WiRouter - Hotspot Manager")
	fmt.Fprintln(output, "========================================")
	fmt.Fprintln(output)
	fmt.Fprintln(output, "Connessioni con accesso a Internet:")
	for index, source := range sources {
		fmt.Fprintf(output, "  %d) %s\n", index+1, source)
	}

	reader := bufio.NewReader(input)
	source, err := readSource(reader, output, sources)
	if err != nil {
		return err
	}
	ssid, err := readValidated(reader, output, "Inserisci SSID hotspot: ", hotspot.ValidateSSID)
	if err != nil {
		return err
	}
	password, err := readValidated(
		reader,
		output,
		"Inserisci password (8-63 caratteri): ",
		hotspot.ValidatePassword,
	)
	if err != nil {
		return err
	}

	// Registrato prima dell'avvio: CTRL+C non deve terminare il processo senza fermare l'hotspot.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	defer signal.Stop(interrupt)

	config := model.HotspotConfig{SSID: ssid, Password: password, Source: source}
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
