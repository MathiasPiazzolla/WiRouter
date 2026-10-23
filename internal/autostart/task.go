package autostart

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"text/template"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

// TaskName è il nome dell'attività pianificata creata da WiRouter.
const TaskName = "WiRouter"

// Options descrive l'attività pianificata da creare.
type Options struct {
	// ExePath è il percorso completo dell'eseguibile da avviare.
	ExePath string
	// UserID è l'utente Windows (DOMINIO\utente) per cui scatta l'avvio.
	UserID string
	// Args sono gli argomenti passati all'eseguibile, ad esempio "--autostart".
	Args string
}

// Template dell'attività: si avvia all'accesso, senza limiti di durata,
// riprova in caso di errore e non si ferma con la batteria.
const taskTemplate = `<?xml version="1.0" encoding="UTF-16"?>
<Task version="1.2" xmlns="http://schemas.microsoft.com/windows/2004/02/mit/task">
  <RegistrationInfo>
    <Description>Avvia WiRouter all'accesso e condivide la connessione Internet tramite hotspot Wi-Fi.</Description>
  </RegistrationInfo>
  <Triggers>
    <LogonTrigger>
      <Enabled>true</Enabled>
      <UserId>{{.UserID}}</UserId>
      <Delay>PT30S</Delay>
    </LogonTrigger>
  </Triggers>
  <Principals>
    <Principal id="Author">
      <UserId>{{.UserID}}</UserId>
      <LogonType>InteractiveToken</LogonType>
      <RunLevel>LeastPrivilege</RunLevel>
    </Principal>
  </Principals>
  <Settings>
    <MultipleInstancesPolicy>IgnoreNew</MultipleInstancesPolicy>
    <DisallowStartIfOnBatteries>false</DisallowStartIfOnBatteries>
    <StopIfGoingOnBatteries>false</StopIfGoingOnBatteries>
    <StartWhenAvailable>true</StartWhenAvailable>
    <RunOnlyIfNetworkAvailable>false</RunOnlyIfNetworkAvailable>
    <AllowStartOnDemand>true</AllowStartOnDemand>
    <Enabled>true</Enabled>
    <Hidden>false</Hidden>
    <ExecutionTimeLimit>PT0S</ExecutionTimeLimit>
    <RestartOnFailure>
      <Interval>PT1M</Interval>
      <Count>3</Count>
    </RestartOnFailure>
  </Settings>
  <Actions Context="Author">
    <Exec>
      <Command>{{.ExePath}}</Command>
      <Arguments>{{.Args}}</Arguments>
    </Exec>
  </Actions>
</Task>
`

// BuildXML restituisce il documento XML dell'attività, con codifica UTF-16 LE
// e BOM, come richiesto da schtasks /xml.
func BuildXML(opt Options) ([]byte, error) {
	tmpl, err := template.New("task").Parse(taskTemplate)
	if err != nil {
		return nil, err
	}
	data := Options{
		ExePath: escapeXML(opt.ExePath),
		UserID:  escapeXML(opt.UserID),
		Args:    escapeXML(opt.Args),
	}
	var doc bytes.Buffer
	if err := tmpl.Execute(&doc, data); err != nil {
		return nil, err
	}
	return encodeUTF16LE(doc.String()), nil
}

func escapeXML(s string) string {
	var b bytes.Buffer
	_ = xml.EscapeText(&b, []byte(s))
	return b.String()
}

func encodeUTF16LE(s string) []byte {
	units := utf16.Encode([]rune(s))
	buf := make([]byte, 2+2*len(units))
	buf[0], buf[1] = 0xFF, 0xFE
	for i, u := range units {
		binary.LittleEndian.PutUint16(buf[2+2*i:], u)
	}
	return buf
}

// Install crea (o sovrascrive) l'attività pianificata per l'utente indicato.
func Install(opt Options) error {
	xmlData, err := BuildXML(opt)
	if err != nil {
		return err
	}
	dir, err := os.MkdirTemp("", "wirouter")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)

	path := filepath.Join(dir, "task.xml")
	if err := os.WriteFile(path, xmlData, 0o600); err != nil {
		return err
	}
	return runSchtasks("/create", "/tn", TaskName, "/xml", path, "/f")
}

// Uninstall rimuove l'attività pianificata di WiRouter.
func Uninstall() error {
	return runSchtasks("/delete", "/tn", TaskName, "/f")
}

// CurrentUserID restituisce DOMINIO\utente dell'utente attuale.
func CurrentUserID() string {
	return fmt.Sprintf(`%s\%s`, os.Getenv("USERDOMAIN"), os.Getenv("USERNAME"))
}

func runSchtasks(args ...string) error {
	cmd := exec.Command("schtasks", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: windows.CREATE_NO_WINDOW,
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("schtasks %v: %w: %s", args, err, string(out))
	}
	return nil
}
