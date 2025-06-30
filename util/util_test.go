package util_test

import (
	"net"
	"reflect"
	"testing"

	"github.com/datum-cloud/galactic/util"
)

func TestParseIP(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantIP    net.IP
		wantError bool
	}{
		{"ValidIPv4", "192.168.0.1", net.ParseIP("192.168.0.1"), false},
		{"ValidIPv6", "2607:ed40:ff00::1", net.ParseIP("2607:ed40:ff00::1"), false},
		{"InvalidIP", "not_an_ip", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.ParseIP(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseIP() error = %v, wantError = %v", err, tt.wantError)
			}
			if !reflect.DeepEqual(got, tt.wantIP) {
				t.Errorf("ParseIP() got = %v, want = %v", got, tt.wantIP)
			}
		})
	}
}

func TestParseSegments(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantIPs   []net.IP
		wantError bool
	}{
		{
			"ValidSingleSegment",
			"2607:ed40:ff00::1",
			[]net.IP{net.ParseIP("2607:ed40:ff00::1")},
			false,
		},
		{
			"ValidMultipleSegments",
			"2607:ed40:ff00::1, 2607:ed40:ff01::1",
			[]net.IP{net.ParseIP("2607:ed40:ff01::1"), net.ParseIP("2607:ed40:ff00::1")},
			false,
		},
		{
			"InvalidSegment",
			"2607:ed40:ff00::1, invalid_ip",
			nil,
			true,
		},
		{
			"InvalidIPv4Segment",
			"2607:ed40:ff00::1, 192.168.0.1",
			nil,
			true,
		},
		{
			"EmptyInput",
			"",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := util.ParseSegments(tt.input)
			if (err != nil) != tt.wantError {
				t.Errorf("ParseSegments() error = %v, wantError = %v", err, tt.wantError)
			}
			if !tt.wantError && !reflect.DeepEqual(got, tt.wantIPs) {
				t.Errorf("ParseSegments() got = %v, want = %v", got, tt.wantIPs)
			}
		})
	}
}

func TestGenerateInterfaceNameVRF(t *testing.T) {
	id := 42
	expected := "galactic42-vrf"
	got := util.GenerateInterfaceNameVRF(id)
	if got != expected {
		t.Errorf("GenerateInterfaceNameVRF(%d) = %s, want %s", id, got, expected)
	}
}

func TestGenerateInterfaceNameHost(t *testing.T) {
	id := 42
	expected := "galactic42-host"
	got := util.GenerateInterfaceNameHost(id)
	if got != expected {
		t.Errorf("GenerateInterfaceNameHost(%d) = %s, want %s", id, got, expected)
	}
}

func TestGenerateInterfaceNameGuest(t *testing.T) {
	id := 42
	expected := "galactic42-guest"
	got := util.GenerateInterfaceNameGuest(id)
	if got != expected {
		t.Errorf("GenerateInterfaceNameGuest(%d) = %s, want %s", id, got, expected)
	}
}
