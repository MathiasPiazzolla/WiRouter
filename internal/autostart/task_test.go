package autostart

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"strings"
	"testing"
	"unicode/utf16"
)

func decodeUTF16LE(t *testing.T, data []byte) string {
	t.Helper()
	if len(data) < 2 || data[0] != 0xFF || data[1] != 0xFE {
		t.Fatalf("manca il BOM UTF-16 LE")
	}
	body := data[2:]
	if len(body)%2 != 0 {
		t.Fatalf("lunghezza UTF-16 dispari")
	}
	units := make([]uint16, len(body)/2)
	for i := range units {
		units[i] = binary.LittleEndian.Uint16(body[2*i:])
	}
	return string(utf16.Decode(units))
}

func TestBuildXMLIsWellFormedAndComplete(t *testing.T) {
	data, err := BuildXML(Options{
		ExePath: `C:\Programmi\WiRouter\wirouter-autostart.exe`,
		UserID:  `DOMINIO\mario`,
		Args:    "--autostart",
	})
	if err != nil {
		t.Fatalf("BuildXML: %v", err)
	}
	doc := decodeUTF16LE(t, data)

	var parsed struct {
		XMLName xml.Name `xml:"Task"`
	}
	if err := xml.Unmarshal([]byte(strings.Replace(doc, `encoding="UTF-16"`, `encoding="UTF-8"`, 1)), &parsed); err != nil {
		t.Fatalf("XML non valido: %v", err)
	}

	checks := []string{
		"<Delay>PT30S</Delay>",
		"<ExecutionTimeLimit>PT0S</ExecutionTimeLimit>",
		"<Interval>PT1M</Interval>",
		"<Count>3</Count>",
		"<RunLevel>LeastPrivilege</RunLevel>",
		"<LogonType>InteractiveToken</LogonType>",
		`<Command>C:\Programmi\WiRouter\wirouter-autostart.exe</Command>`,
		"<Arguments>--autostart</Arguments>",
		`<UserId>DOMINIO\mario</UserId>`,
	}
	for _, want := range checks {
		if !strings.Contains(doc, want) {
			t.Errorf("XML senza %q", want)
		}
	}
}

func TestBuildXMLEscapesValues(t *testing.T) {
	data, err := BuildXML(Options{
		ExePath: `C:\A&B\wirouter.exe`,
		UserID:  `D\u<x>`,
		Args:    `--a "b" & c`,
	})
	if err != nil {
		t.Fatalf("BuildXML: %v", err)
	}
	doc := decodeUTF16LE(t, data)
	if !bytes.Contains([]byte(doc), []byte("C:\\A&amp;B\\wirouter.exe")) {
		t.Errorf("percorso non escapato correttamente: %s", doc)
	}
	var parsed struct{}
	if err := xml.Unmarshal([]byte(strings.Replace(doc, `encoding="UTF-16"`, `encoding="UTF-8"`, 1)), &parsed); err != nil {
		t.Fatalf("XML con valori speciali non valido: %v", err)
	}
}
