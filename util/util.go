package util

import (
	"fmt"
	"net"
)

const InterfaceNameTemplate = "galactic%d-%s"

func ParseIP(ip string) (net.IP, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, fmt.Errorf("cannot parse IP: %v", ip)
	}
	return parsed, nil
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
