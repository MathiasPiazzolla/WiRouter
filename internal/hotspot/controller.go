package hotspot

import "github.com/MathiasPiazzolla/WiRouter/model"

// Controller adatta le funzioni del pacchetto all'interfaccia usata dal runner.
type Controller struct{}

func (Controller) ListSources() ([]string, error) {
	return ListSources()
}

func (Controller) Status(source string) (model.HotspotStatus, error) {
	return GetHotspotStatus(source)
}

func (Controller) Start(config model.HotspotConfig) error {
	return StartHotspot(config)
}

func (Controller) Stop(source string) error {
	return StopHotspot(source)
}
