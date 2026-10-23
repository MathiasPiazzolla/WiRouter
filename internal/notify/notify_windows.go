package notify

import "golang.org/x/sys/windows"

// Error mostra una finestra di errore di Windows. Serve perché all'avvio
// automatico il programma non ha una console visibile.
func Error(title, text string) error {
	title16, err := windows.UTF16PtrFromString(title)
	if err != nil {
		return err
	}
	text16, err := windows.UTF16PtrFromString(text)
	if err != nil {
		return err
	}
	ret, err := windows.MessageBox(0, text16, title16, windows.MB_OK|windows.MB_ICONERROR|windows.MB_TOPMOST)
	if ret == 0 {
		return err
	}
	return nil
}
