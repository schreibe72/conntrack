package main

import (
	"errors"
	"net"
)

type Interface struct {
	iface net.Interface
	ips   []net.IP
}

func getLocalInterface(device string) (Interface, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return Interface{}, err
	}
	for _, i := range interfaces {
		if i.Name == device {
			result := Interface{
				iface: i,
			}
			addrs, err := i.Addrs()
			if err != nil {
				return Interface{}, err
			}
			for _, a := range addrs {
				ip, _, _ := net.ParseCIDR(a.String())
				result.ips = append(result.ips, ip)
			}
			return result, nil
		}
	}
	return Interface{}, errors.New("No Device Found")
}

func (i *Interface) isLocalIP(ip net.IP) bool {
	for _, localip := range i.ips {
		if localip.Equal(ip) {
			return true
		}
	}
	return false
}
