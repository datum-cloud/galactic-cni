package util

import (
	"fmt"
	"net"
	"strings"
)

const InterfaceNameTemplate = "galactic%d-%s"

func ParseIP(ip string) (net.IP, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, fmt.Errorf("cannot parse IP: %v", ip)
	}
	return parsed, nil
}

func ParseSegments(input string) ([]net.IP, error) {
	var segments []net.IP
	for _, s := range strings.Split(input, ",") {
		s = strings.TrimSpace(s)
		ip, err := ParseIP(s)
		if err != nil {
			return nil, fmt.Errorf("could not parse ip (%s): %v", s, err)
		}
		segments = append([]net.IP{ip}, segments...)
	}
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments parsed: %v", input)
	}
	return segments, nil
}

func GenerateInterfaceNameVRF(id int) string {
	return fmt.Sprintf(InterfaceNameTemplate, id, "vrf")
}

func GenerateInterfaceNameHost(id int) string {
	return fmt.Sprintf(InterfaceNameTemplate, id, "host")
}

func GenerateInterfaceNameGuest(id int) string {
	return fmt.Sprintf(InterfaceNameTemplate, id, "guest")
}
