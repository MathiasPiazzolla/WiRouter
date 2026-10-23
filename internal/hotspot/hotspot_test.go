package hotspot

import (
	"errors"
	"strings"
	"testing"

	"github.com/MathiasPiazzolla/WiRouter/model"
)

func stubExecuteCommand(t *testing.T, stub func(env []string, name string, args ...string) ([]byte, error)) {
	t.Helper()
	original := executeCommand
	executeCommand = stub
	t.Cleanup(func() { executeCommand = original })
}

func TestValidateSSID(t *testing.T) {
	tests := []struct {
		name    string
		ssid    string
		wantErr bool
	}{
		{name: "valid", ssid: "WiRouter"},
		{name: "empty", ssid: "   ", wantErr: true},
		{name: "too long", ssid: strings.Repeat("a", 33), wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := ValidateSSID(test.ssid); (err != nil) != test.wantErr {
				t.Fatalf("ValidateSSID(%q) error = %v, wantErr %v", test.ssid, err, test.wantErr)
			}
		})
	}
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{name: "valid", password: "password123"},
		{name: "too short", password: "short", wantErr: true},
		{name: "too long", password: strings.Repeat("a", 64), wantErr: true},
		{name: "non ASCII", password: "passwordè", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := ValidatePassword(test.password); (err != nil) != test.wantErr {
				t.Fatalf("ValidatePassword(%q) error = %v, wantErr %v", test.password, err, test.wantErr)
			}
		})
	}
}

func TestStartHotspotPassesInputsOnlyThroughEnvironment(t *testing.T) {
	var gotEnv []string
	var gotName string
	var gotScript string
	stubExecuteCommand(t, func(env []string, name string, args ...string) ([]byte, error) {
		gotEnv, gotName, gotScript = env, name, args[len(args)-1]
		return nil, nil
	})

	ssid := `Wi'Router"; Remove-Item *`
	err := StartHotspot(model.HotspotConfig{SSID: ssid, Password: "password123", Source: "Ethernet"})
	if err != nil {
		t.Fatalf("StartHotspot() error = %v", err)
	}

	if gotName != "powershell.exe" {
		t.Errorf("command = %q, want powershell.exe", gotName)
	}
	wantEnv := []string{"WIROUTER_SOURCE=Ethernet", "WIROUTER_SSID=" + ssid, "WIROUTER_PASSWORD=password123"}
	if strings.Join(gotEnv, "|") != strings.Join(wantEnv, "|") {
		t.Errorf("env = %#v, want %#v", gotEnv, wantEnv)
	}
	if strings.Contains(gotScript, ssid) || strings.Contains(gotScript, "password123") {
		t.Error("lo script contiene input utente")
	}
	if !strings.Contains(gotScript, "StartTetheringAsync") {
		t.Error("lo script non avvia il tethering")
	}
}

func TestStartHotspotRejectsInvalidInputWithoutRunning(t *testing.T) {
	stubExecuteCommand(t, func([]string, string, ...string) ([]byte, error) {
		t.Fatal("il comando non deve essere eseguito")
		return nil, nil
	})

	if err := StartHotspot(model.HotspotConfig{SSID: "", Password: "password123"}); err == nil {
		t.Error("SSID vuoto accettato")
	}
	if err := StartHotspot(model.HotspotConfig{SSID: "WiRouter", Password: "short"}); err == nil {
		t.Error("password corta accettata")
	}
}

func TestRunScriptReturnsScriptErrorMessage(t *testing.T) {
	stubExecuteCommand(t, func([]string, string, ...string) ([]byte, error) {
		return []byte("Hotspot già attivo\n"), errors.New("exit status 1")
	})

	err := StopHotspot("")
	if err == nil || !strings.Contains(err.Error(), "Hotspot già attivo") {
		t.Fatalf("StopHotspot() error = %v, want script message", err)
	}
}

func TestListSources(t *testing.T) {
	stubExecuteCommand(t, func([]string, string, ...string) ([]byte, error) {
		return []byte("Ethernet\r\n\r\nWi-Fi casa\r\n"), nil
	})

	sources, err := ListSources()
	if err != nil {
		t.Fatalf("ListSources() error = %v", err)
	}
	if len(sources) != 2 || sources[0] != "Ethernet" || sources[1] != "Wi-Fi casa" {
		t.Fatalf("ListSources() = %#v", sources)
	}
}

func TestParseStatus(t *testing.T) {
	status := parseStatus("state=On\r\nssid=WiRouter\r\nclients=3\r\n")

	if !status.IsRunning || status.SSID != "WiRouter" || status.ConnectedDevices != 3 {
		t.Fatalf("parseStatus() = %+v", status)
	}
}
