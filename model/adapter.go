package model

import (
	"net"
)

type Adapter struct {
	Name         string
	HardwareAddr string
}

type AdapterWin struct {
	InterfaceID uint32
	Name        string
	Description string
	MAC         string
	Type        string
}

func GetAllAdapters() ([]Adapter, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var adapters []Adapter
	for _, iface := range interfaces {
		adapters = append(adapters, Adapter{
			Name:         iface.Name,
			HardwareAddr: iface.HardwareAddr.String(),
		})
	}
	return adapters, nil
}
