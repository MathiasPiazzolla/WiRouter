package main

import (
	"reflect"
	"testing"

	"golang.org/x/sys/windows"

	"github.com/MathiasPiazzolla/WiRouter/model"
)

func TestAutostartArgsSurviveCommandLineParsing(t *testing.T) {
	config := model.HotspotConfig{
		Source:   "Connessione alla rete locale",
		SSID:     `Casa "Rossi" \`,
		Password: `pa ss"word\`,
	}

	line := "wirouter-autostart.exe " + autostartArgs(config)
	got, err := windows.DecomposeCommandLine(line)
	if err != nil {
		t.Fatalf("DecomposeCommandLine: %v", err)
	}

	want := []string{
		"wirouter-autostart.exe",
		"-autostart",
		"-source", config.Source,
		"-ssid", config.SSID,
		"-password", config.Password,
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("argomenti: got %q, want %q", got, want)
	}
}
