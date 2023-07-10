package main

import (
	"eaproxy/pkg/config"
	"eaproxy/pkg/eaproxy"
	"eaproxy/pkg/socket"
	"fmt"
	"os"
)

var eapMulticastAddr = []uint8{0x01, 0x80, 0xc2, 0x00, 0x00, 0x03}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s [path to config toml file]\n", os.Args[0])
		os.Exit(1)
	}
	config, err := config.LoadConfig(os.Args[1])
	if err != nil {
		panic(err)
	}
	if config.WanIfname == "" {
		fmt.Printf("configuration must have a 'wan_ifname' property\n")
		os.Exit(1)
	}
	if config.RouterIfname == "" {
		fmt.Printf("configuration must have a 'router_ifname' property\n")
		os.Exit(1)
	}
	if config.VlanID < 0 {
		fmt.Printf("configuration must have a non-negative 'vlan_id' property\n")
		os.Exit(1)
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
