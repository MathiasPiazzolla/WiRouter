package instance

import (
	"errors"
	"testing"
)

func TestAcquireSecondInstanceFails(t *testing.T) {
	name := `Local\WiRouterTest-` + t.Name()

	release, err := Acquire(name)
	if err != nil {
		t.Fatalf("prima istanza: %v", err)
	}
	defer release()

	releaseSecond, err := Acquire(name)
	if !errors.Is(err, ErrAlreadyRunning) {
		if releaseSecond != nil {
			releaseSecond()
		}
		t.Fatalf("seconda istanza: errore %v, atteso ErrAlreadyRunning", err)
	}
}
