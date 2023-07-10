package mac

import (
	"eaproxy/pkg/vyatta"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

// Get gets a mac address of a specific interface
func Get(ifname string) ([]byte, error) {
	filename := fmt.Sprintf("/sys/class/net/%s/address", ifname)
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	str := string(buf)
	octetStrs := strings.Split(str, ":")
	var octets []byte
	for i, octetStr := range octetStrs {
		ui64, err := strconv.ParseUint(octetStr, 16, 8)
		if err != nil {
			return nil, err
		}
		ui8 := uint8(ui64)
		octets[i] = ui8
	}
	return octets, nil
}

// Set sets the mac address of a specific interface
func Set(ifname string, mac string) error {
	return vyatta.Run(ifname, "--set-mac", mac)
}
