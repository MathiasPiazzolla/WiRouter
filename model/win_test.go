package model

import "testing"

func TestFormatMAC(t *testing.T) {
	t.Parallel()

	if got := formatMAC([]byte{0x00, 0x1A, 0x2B, 0x03}); got != "00:1A:2B:03" {
		t.Fatalf("formatMAC() = %q, want %q", got, "00:1A:2B:03")
	}
	if got := formatMAC(nil); got != "" {
		t.Fatalf("formatMAC(nil) = %q, want empty string", got)
	}
}

func TestInterfaceType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		value uint32
		want  string
	}{
		{value: 6, want: "Ethernet"},
		{value: 71, want: "Wi-Fi"},
		{value: 24, want: "Loopback"},
		{value: 999, want: "Unknown (999)"},
	}

	for _, test := range tests {
		if got := interfaceType(test.value); got != test.want {
			t.Errorf("interfaceType(%d) = %q, want %q", test.value, got, test.want)
		}
	}
}
