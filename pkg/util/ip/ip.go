package ip

import (
	"net"
)

func LocalIp() (i net.IP, err error) {
	addrs, err := net.InterfaceAddrs()
	for _, v := range addrs {
		if ipnet, ok := v.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				i = ipnet.IP
				break
			}
		}
	}
	return
}
