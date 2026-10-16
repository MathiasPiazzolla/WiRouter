package model

import (
	"errors"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

const maxAdapterQueryAttempts = 3

func interfaceType(value uint32) string {
	switch value {
	case windows.IF_TYPE_ETHERNET_CSMACD:
		return "Ethernet"
	case windows.IF_TYPE_IEEE80211:
		return "Wi-Fi"
	case windows.IF_TYPE_SOFTWARE_LOOPBACK:
		return "Loopback"
	default:
		return fmt.Sprintf("Unknown (%d)", value)
	}
}

func formatMAC(address []byte) string {
	if len(address) == 0 {
		return ""
	}

	result := fmt.Sprintf("%02X", address[0])
	for _, value := range address[1:] {
		result += fmt.Sprintf(":%02X", value)
	}
	return result
}

func ListWindowsAdapters() ([]AdapterWin, error) {
	var size uint32
	var flags uint32
	err := windows.GetAdaptersAddresses(windows.AF_UNSPEC, flags, 0, nil, &size)
	if !errors.Is(err, windows.ERROR_BUFFER_OVERFLOW) {
		if err == nil {
			return nil, errors.New("GetAdaptersAddresses non ha restituito la dimensione richiesta")
		}
		return nil, fmt.Errorf("lettura dimensione adapter: %w", err)
	}

	for attempt := 0; attempt < maxAdapterQueryAttempts; attempt++ {
		if size == 0 {
			return nil, errors.New("dimensione buffer adapter non valida")
		}

		buffer := make([]byte, size)
		first := (*windows.IpAdapterAddresses)(unsafe.Pointer(&buffer[0]))
		err = windows.GetAdaptersAddresses(windows.AF_UNSPEC, flags, 0, first, &size)
		if errors.Is(err, windows.ERROR_BUFFER_OVERFLOW) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("lettura adapter: %w", err)
		}

		var adapters []AdapterWin
		for current := first; current != nil; current = current.Next {
			physicalAddressLength := int(current.PhysicalAddressLength)
			if physicalAddressLength > len(current.PhysicalAddress) {
				physicalAddressLength = len(current.PhysicalAddress)
			}

			adapters = append(adapters, AdapterWin{
				InterfaceID: current.IfIndex,
				Name:        utf16String(current.FriendlyName),
				Description: utf16String(current.Description),
				MAC:         formatMAC(current.PhysicalAddress[:physicalAddressLength]),
				Type:        interfaceType(current.IfType),
			})
		}
		return adapters, nil
	}

	return nil, errors.New("la lista adapter è cambiata durante la lettura")
}

func utf16String(value *uint16) string {
	if value == nil {
		return ""
	}
	return windows.UTF16PtrToString(value)
}
