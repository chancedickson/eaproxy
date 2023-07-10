package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Config struct
type Config struct {
	WanIfname    string `toml:"wan_ifname"`
	RouterIfname string `toml:"router_ifname"`
	VlanID       int    `toml:"vlan_id" default:"-1"`
}

// LoadConfig reads a file and parses it as toml into an eaproxy.Config
func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := Config{}
	toml.Unmarshal(data, &config)
	return &config, nil
}
