package utils

import (
	"fmt"
	"net"
)

func ParseMAC(macStr string) (net.HardwareAddr, error) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат MAC-адреса: %s", macStr)
	}
	return mac, nil
}

func ParseIP(ipStr string) (net.IP, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, fmt.Errorf("неверный формат IP-адреса: %s", ipStr)
	}
	return ip, nil
}
