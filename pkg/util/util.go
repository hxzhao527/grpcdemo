package util

import (
	"net"
	"path"
	"strings"
)

func GetServiceNameFromFullMethod(fm string) string {
	return strings.Trim(path.Dir(fm), "/")
}

// GetSelfIPAddress maybe not work well.
func GetSelfIPAddress() net.IP {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				if !v.IP.IsLoopback() {
					return v.IP
				}
			case *net.IPAddr:
				if !v.IP.IsLoopback() {
					return v.IP
				}
			}
		}
	}
	return nil
}
