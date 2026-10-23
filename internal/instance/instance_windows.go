package instance

import (
	"errors"

	"golang.org/x/sys/windows"
)

// ErrAlreadyRunning indica che un'altra istanza di WiRouter è già in esecuzione.
var ErrAlreadyRunning = errors.New("WiRouter è già in esecuzione")

// Acquire crea un mutex con nome e restituisce la funzione che lo rilascia.
// Il mutex va tenuto per tutta la vita del processo.
func Acquire(name string) (release func(), err error) {
	name16, err := windows.UTF16PtrFromString(name)
	if err != nil {
		return nil, err
	}
	handle, err := windows.CreateMutex(nil, false, name16)
	if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
		_ = windows.CloseHandle(handle)
		return nil, ErrAlreadyRunning
	}
	if err != nil {
		return nil, err
	}
	return func() { _ = windows.CloseHandle(handle) }, nil
}
