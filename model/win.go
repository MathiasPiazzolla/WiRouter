package model

import (
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	iphlpapi = windows.NewLazySystemDLL("iphlpapi.dll")

	getAdaptersAddresses = iphlpapi.NewProc("GetAdaptersAddresses")
)

type ipAdapterAddresses struct {
	Length  uint32
	IfIndex uint32

	Next uintptr

	AdapterName uintptr

	FirstUnicastAddress   uintptr
	FirstAnycastAddress   uintptr
	FirstMulticastAddress uintptr
	FirstDnsServerAddress uintptr

	DnsSuffix    uintptr
	Description  uintptr
	FriendlyName uintptr

	PhysicalAddress       [8]byte
	PhysicalAddressLength uint32

	Flags uint32

	Mtu    uint32
	IfType uint32
}

func interfaceType(t uint32) string {

	switch t {

	case 6:
		return "Ethernet"

	case 71:
		return "Wi-Fi"

	case 24:
		return "Loopback"

	default:
		return fmt.Sprintf("Unknown (%d)", t)
	}
}

func formatMAC(mac [8]byte, length uint32) string {

	result := ""

	for i := uint32(0); i < length; i++ {

		if i > 0 {
			result += ":"
		}

		result += fmt.Sprintf("%02X", mac[i])
	}

	return result
}

func ListWindowsAdapters() ([]AdapterWin, error) {

	var size uint32

	// prima chiamata per sapere la dimensione
	getAdaptersAddresses.Call(0, 0, 0, 0, uintptr(unsafe.Pointer(&size)))

	buffer := make([]byte, size)

	ret, _, err := getAdaptersAddresses.Call(0, 0, 0, uintptr(unsafe.Pointer(&buffer[0])), uintptr(unsafe.Pointer(&size)))

	if ret != 0 {
		return nil, err
	}

	first := (*ipAdapterAddresses)(
		unsafe.Pointer(&buffer[0]),
	)

	var result []AdapterWin

	current := first

	for current != nil {

		name := windows.UTF16PtrToString(
			(*uint16)(unsafe.Pointer(current.FriendlyName)),
		)

		description := windows.UTF16PtrToString(
			(*uint16)(unsafe.Pointer(current.Description)),
		)

		result = append(result, AdapterWin{
			InterfaceID: current.IfIndex,
			Name:        name,
			Description: description,
			MAC:         formatMAC(current.PhysicalAddress, current.PhysicalAddressLength),
			Type:        interfaceType(current.IfType),
		})

		if current.Next == 0 {
			break
		}

		current = (*ipAdapterAddresses)(
			unsafe.Pointer(current.Next),
		)
	}

	return result, nil

}
