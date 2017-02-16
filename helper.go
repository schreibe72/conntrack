package main

import (
	"net"
	"os"
	"path/filepath"
)

func isInternet(ip net.IP) bool {
	_, subnet192, _ := net.ParseCIDR("192.168.0.0/16")
	_, subnet172, _ := net.ParseCIDR("172.16.0.0/12")
	_, subnet10, _ := net.ParseCIDR("10.0.0.0/8")
	if subnet192.Contains(ip) {
		return false
	}
	if subnet172.Contains(ip) {
		return false
	}
	if subnet10.Contains(ip) {
		return false
	}
	return true
}

func uint16InSlice(list []uint16, port uint16) bool {
	for _, b := range list {
		if b == port {
			return true
		}
	}
	return false
}

func createDataPath(path string) {
	mode := os.FileMode(0755)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, mode)
	}
	directions := [3]string{"in", "out", "else"}
	protos := [2]string{"tcp", "udp"}
	for _, dir1 := range directions {
		path1 := filepath.Join(path, dir1)
		if _, err := os.Stat(path1); os.IsNotExist(err) {
			os.Mkdir(path1, mode)
		}
		for _, dir2 := range protos {
			path2 := filepath.Join(path1, dir2)
			if _, err := os.Stat(path2); os.IsNotExist(err) {
				os.Mkdir(path2, mode)
			}
		}
	}
}
