package hotspot

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/MathiasPiazzolla/WiRouter/model"
)

const (
	minPasswordLength = 8
	maxPasswordLength = 63
	maxSSIDBytes      = 32
)

var executeCommand = func(env []string, name string, args ...string) ([]byte, error) {
	command := exec.Command(name, args...)
	command.Env = append(os.Environ(), env...)
	return command.CombinedOutput()
}

func ValidateSSID(ssid string) error {
	trimmed := strings.TrimSpace(ssid)
	if trimmed == "" {
		return errors.New("l'SSID non può essere vuoto")
	}
	if len(trimmed) > maxSSIDBytes {
		return fmt.Errorf("l'SSID non può superare %d byte", maxSSIDBytes)
	}
	return nil
}

func ValidatePassword(password string) error {
	length := len(password)
	if length < minPasswordLength {
		return fmt.Errorf("la password deve contenere almeno %d caratteri", minPasswordLength)
	}
	if length > maxPasswordLength {
		return fmt.Errorf("la password non può superare %d caratteri", maxPasswordLength)
	}
	for _, character := range password {
		if character < 32 || character > 126 {
			return errors.New("la password deve contenere solo caratteri ASCII stampabili")
		}
	}
	return nil
}

// ListSources restituisce i nomi delle connessioni con accesso a Internet.
func ListSources() ([]string, error) {
	output, err := runScript(nil, scriptSources)
	if err != nil {
		return nil, fmt.Errorf("lettura connessioni Internet: %w", err)
	}

	var sources []string
	for _, line := range strings.Split(output, "\n") {
		if name := strings.TrimSpace(line); name != "" {
			sources = append(sources, name)
		}
	}
	return sources, nil
}

// StartHotspot configura e avvia l'hotspot condividendo la connessione config.Source.
func StartHotspot(config model.HotspotConfig) error {
	if err := ValidateSSID(config.SSID); err != nil {
		return err
	}
	if err := ValidatePassword(config.Password); err != nil {
		return err
	}

	env := []string{
		"WIROUTER_SOURCE=" + config.Source,
		"WIROUTER_SSID=" + strings.TrimSpace(config.SSID),
		"WIROUTER_PASSWORD=" + config.Password,
	}
	if _, err := runScript(env, scriptStart); err != nil {
		return fmt.Errorf("avvio hotspot: %w", err)
	}
	return nil
}

func StopHotspot(source string) error {
	if _, err := runScript(sourceEnv(source), scriptStop); err != nil {
		return fmt.Errorf("arresto hotspot: %w", err)
	}
	return nil
}

func GetHotspotStatus(source string) (model.HotspotStatus, error) {
	output, err := runScript(sourceEnv(source), scriptStatus)
	if err != nil {
		return model.HotspotStatus{}, fmt.Errorf("lettura stato hotspot: %w", err)
	}
	return parseStatus(output), nil
}

func sourceEnv(source string) []string {
	return []string{"WIROUTER_SOURCE=" + source}
}

func runScript(env []string, body string) (string, error) {
	output, err := executeCommand(
		env,
		"powershell.exe",
		"-NoProfile",
		"-NonInteractive",
		"-Command",
		buildScript(body),
	)
	text := strings.TrimSpace(string(output))
	if err != nil {
		if text == "" {
			return "", err
		}
		return "", errors.New(text)
	}
	return text, nil
}

func parseStatus(output string) model.HotspotStatus {
	var status model.HotspotStatus
	for _, line := range strings.Split(output, "\n") {
		key, value, found := strings.Cut(strings.TrimSpace(line), "=")
		if !found {
			continue
		}

		switch key {
		case "state":
			status.IsRunning = strings.EqualFold(value, "On")
		case "ssid":
			status.SSID = value
		case "clients":
			status.ConnectedDevices, _ = strconv.Atoi(value)
		}
	}
	return status
}
