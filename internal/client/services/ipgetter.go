package services

import (
	"fmt"
	"net"
	"sync"
)

type IPGetter struct {
	IP  string
	Err error
}

var instance *IPGetter
var once sync.Once

func GetIpGetter() *IPGetter {
	once.Do(func() {
		instance = &IPGetter{}
		instance.IP, instance.Err = getLocalIP()
	})
	return instance
}

func getLocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", fmt.Errorf("не удалось установить соединение: %v", err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}
