package config

import (
	"io/ioutil"

	"github.com/pelletier/go-toml"
)

// Config struct
type Config struct {
	WanIfname    string
	RouterIfname string
	VlanID       int
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
