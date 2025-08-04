package utils

import (
	"bytes"
	"net"
)

// возвращает список айпи для устройства
func GetInterfaceIPs(name string) ([]net.IP, error) {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return nil, err
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return nil, err
	}

	var ips []net.IP
	for _, addr := range addrs {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}
		// Фильтруем только IPv4-адреса
		if ip.To4() != nil {
			ips = append(ips, ip.To4())
		}
	}
	return ips, nil
}

// проверяет айпи
func IsMyIP(target [4]byte, myIPs []net.IP) bool {
	for _, ip := range myIPs {
		if len(ip) == 4 && ip[0] == target[0] && ip[1] == target[1] && ip[2] == target[2] && ip[3] == target[3] {
			return true
		}
	}
	return false
}

// проверяет является ли мак броудкастом
func IsBroadcast(mac []byte) bool {
	return bytes.Equal(mac, []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff})
}
