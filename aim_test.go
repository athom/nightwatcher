package nightwatcher

import (
	"testing"

	"net"
)

func TestAim_FromIP(t *testing.T) {
	aim := &Aim{}
	localIp := aim.FromIP()

	found := false
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if localIp == ipnet.IP.String() {
					found = true
				}
				println(ipnet.IP.String())
			}
		}
	}
	if !found {
		t.Fatalf("FromIp failed, got %v", localIp)
	}
}
