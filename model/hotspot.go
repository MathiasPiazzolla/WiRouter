package model

type HotspotConfig struct {
	SSID     string
	Password string
	Source   string // nome della connessione Internet da condividere
}

type HotspotStatus struct {
	IsRunning        bool
	SSID             string
	ConnectedDevices int
}
