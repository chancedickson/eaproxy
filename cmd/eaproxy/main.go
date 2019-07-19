package main

import (
	"eaproxy/pkg/config"
	"eaproxy/pkg/eaproxy"
	"eaproxy/pkg/socket"
	"os"
)

var eapMulticastAddr = []uint8{0x01, 0x80, 0xc2, 0x00, 0x00, 0x03}

func main() {
	config, err := config.LoadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}

	routersock, err := socket.New(config.RouterIfname, eapMulticastAddr)
	if err != nil {
		panic(err)
	}
	defer routersock.Close()

	wansock, err := socket.New(config.WanIfname, eapMulticastAddr)
	if err != nil {
		panic(err)
	}
	defer wansock.Close()

	proxy := eaproxy.New(routersock, wansock, config.VlanID)
	err = proxy.Start()
	if err != nil {
		panic(err)
	}
}
