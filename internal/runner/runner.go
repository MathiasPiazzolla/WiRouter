package runner

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MathiasPiazzolla/WiRouter/model"
)

// Hotspot è l'insieme delle operazioni sull'hotspot usate dal runner.
type Hotspot interface {
	ListSources() ([]string, error)
	Status(source string) (model.HotspotStatus, error)
	Start(config model.HotspotConfig) error
	Stop(source string) error
}

// Options regola i tempi e i tentativi del runner.
type Options struct {
	Config          model.HotspotConfig
	SourceTimeout   time.Duration // quanto attendere che la fonte compaia
	SourcePoll      time.Duration // ogni quanto controllare le fonti
	StartAttempts   int           // numero massimo di tentativi di avvio
	StartRetryDelay time.Duration // attesa tra un tentativo e l'altro
}

// Run attende la fonte, riavvia l'hotspot se era già acceso, lo avvia e resta
// in attesa finché ctx non viene annullato. Allo stop ferma l'hotspot.
func Run(ctx context.Context, h Hotspot, opt Options) error {
	source := opt.Config.Source

	if err := waitForSource(ctx, h, source, opt); err != nil {
		return err
	}

	status, err := h.Status(source)
	if err != nil {
		return fmt.Errorf("stato hotspot: %w", err)
	}
	if status.IsRunning {
		if err := h.Stop(source); err != nil {
			return fmt.Errorf("arresto hotspot già attivo: %w", err)
		}
	}

	if err := startWithRetries(ctx, h, opt); err != nil {
		return err
	}

	<-ctx.Done()
	if err := h.Stop(source); err != nil {
		return fmt.Errorf("arresto hotspot: %w", err)
	}
	return nil
}

func waitForSource(ctx context.Context, h Hotspot, source string, opt Options) error {
	deadline := time.NewTimer(opt.SourceTimeout)
	defer deadline.Stop()
	poll := time.NewTicker(opt.SourcePoll)
	defer poll.Stop()

	for {
		sources, err := h.ListSources()
		if err == nil && contains(sources, source) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-deadline.C:
			return fmt.Errorf("fonte %q non disponibile entro %s", source, opt.SourceTimeout)
		case <-poll.C:
		}
	}
}

func startWithRetries(ctx context.Context, h Hotspot, opt Options) error {
	var lastErr error
	for attempt := 1; attempt <= opt.StartAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return bestEffortStop(h, opt, err)
		}
		lastErr = h.Start(opt.Config)
		if lastErr == nil {
			return nil
		}
		if attempt < opt.StartAttempts {
			select {
			case <-ctx.Done():
				return bestEffortStop(h, opt, ctx.Err())
			case <-time.After(opt.StartRetryDelay):
			}
		}
	}
	return bestEffortStop(h, opt, fmt.Errorf("avvio hotspot fallito dopo %d tentativi: %w", opt.StartAttempts, lastErr))
}

// bestEffortStop prova a fermare l'hotspot dopo un errore: il tentativo di avvio
// può essere riuscito solo in parte. Restituisce l'errore originale.
func bestEffortStop(h Hotspot, opt Options, cause error) error {
	if stopErr := h.Stop(opt.Config.Source); stopErr != nil {
		return errors.Join(cause, fmt.Errorf("arresto dopo errore: %w", stopErr))
	}
	return cause
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
