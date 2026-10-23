package runner

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/MathiasPiazzolla/WiRouter/model"
)

type fakeHotspot struct {
	mu            sync.Mutex
	sourceAfter   int  // numero di controlli prima che la fonte compaia (-1 = mai)
	running       bool // stato iniziale dell'hotspot
	startFailures int  // quanti Start devono fallire prima di riuscire
	listCalls     int
	starts        int
	stops         []string
	started       chan struct{}
}

func (f *fakeHotspot) ListSources() ([]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.listCalls++
	if f.sourceAfter >= 0 && f.listCalls > f.sourceAfter {
		return []string{"Ethernet", "Wi-Fi"}, nil
	}
	return []string{"Wi-Fi"}, nil
}

func (f *fakeHotspot) Status(source string) (model.HotspotStatus, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	return model.HotspotStatus{IsRunning: f.running}, nil
}

func (f *fakeHotspot) Start(model.HotspotConfig) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.starts++
	if f.starts <= f.startFailures {
		return errors.New("avvio non riuscito")
	}
	f.running = true
	if f.started != nil {
		select {
		case <-f.started:
		default:
			close(f.started)
		}
	}
	return nil
}

func (f *fakeHotspot) Stop(source string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.stops = append(f.stops, source)
	f.running = false
	return nil
}

func baseOptions(source string) Options {
	return Options{
		Config:          model.HotspotConfig{SSID: "Test", Password: "password1", Source: source},
		SourceTimeout:   2 * time.Second,
		SourcePoll:      time.Millisecond,
		StartAttempts:   3,
		StartRetryDelay: time.Millisecond,
	}
}

// runUntilStarted esegue Run e lo ferma dopo l'avvio dell'hotspot.
func runUntilStarted(t *testing.T, h *fakeHotspot, opt Options) error {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	h.started = make(chan struct{})

	done := make(chan error, 1)
	go func() { done <- Run(ctx, h, opt) }()

	select {
	case <-h.started:
		cancel()
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: l'hotspot non è stato avviato")
	}
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("timeout: Run non è terminato dopo l'annullamento")
		return nil
	}
}

func TestRunStartsWhenSourceAppearsAndStopsOnCancel(t *testing.T) {
	h := &fakeHotspot{sourceAfter: 3}
	err := runUntilStarted(t, h, baseOptions("Ethernet"))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if h.starts != 1 {
		t.Errorf("avvii = %d, atteso 1", h.starts)
	}
	if len(h.stops) != 1 || h.stops[0] != "Ethernet" {
		t.Errorf("stop = %v, atteso [Ethernet]", h.stops)
	}
}

func TestRunTimesOutWhenSourceNeverAppears(t *testing.T) {
	h := &fakeHotspot{sourceAfter: -1}
	opt := baseOptions("Ethernet")
	opt.SourceTimeout = 20 * time.Millisecond

	err := Run(context.Background(), h, opt)
	if err == nil || !strings.Contains(err.Error(), "non disponibile") {
		t.Fatalf("errore = %v, atteso timeout della fonte", err)
	}
	if h.starts != 0 {
		t.Errorf("avvii = %d, attesi 0", h.starts)
	}
}

func TestRunStopsRunningHotspotBeforeStart(t *testing.T) {
	h := &fakeHotspot{sourceAfter: 0, running: true}
	err := runUntilStarted(t, h, baseOptions("Ethernet"))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if h.starts != 1 {
		t.Errorf("avvii = %d, atteso 1", h.starts)
	}
	if len(h.stops) != 2 {
		t.Errorf("stop = %v, attesi due (prima dell'avvio e alla chiusura)", h.stops)
	}
}

func TestRunRetriesStartUntilItWorks(t *testing.T) {
	h := &fakeHotspot{sourceAfter: 0, startFailures: 2}
	err := runUntilStarted(t, h, baseOptions("Ethernet"))
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if h.starts != 3 {
		t.Errorf("avvii = %d, attesi 3", h.starts)
	}
}

func TestRunFailsAfterAllStartAttempts(t *testing.T) {
	h := &fakeHotspot{sourceAfter: 0, startFailures: 100}
	err := Run(context.Background(), h, baseOptions("Ethernet"))
	if err == nil || !strings.Contains(err.Error(), "3 tentativi") {
		t.Fatalf("errore = %v, atteso fallimento dopo 3 tentativi", err)
	}
	if h.starts != 3 {
		t.Errorf("avvii = %d, attesi 3", h.starts)
	}
	if len(h.stops) != 1 {
		t.Errorf("stop = %v, atteso un arresto di sicurezza", h.stops)
	}
}

func TestRunCancelledWhileWaitingForSource(t *testing.T) {
	h := &fakeHotspot{sourceAfter: -1}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Run(ctx, h, baseOptions("Ethernet"))
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("errore = %v, atteso context.Canceled", err)
	}
	if h.starts != 0 || len(h.stops) != 0 {
		t.Errorf("avvii = %d, stop = %v: attesi nessun avvio e nessun arresto", h.starts, h.stops)
	}
}
