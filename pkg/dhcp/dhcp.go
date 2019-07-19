package dhcp

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/sys/unix"
)

func dhclientPathnames(ifname string) (string, string, string) {
	f := strings.Replace(ifname, ".", "_", -1)
	return fmt.Sprintf("/var/run/dhclient_%s.conf", f),
		fmt.Sprintf("/var/run/dhclient_%s.pid", f),
		fmt.Sprintf("/var/run/dhclient_%s.leases", f)
}

// Obtain obtains a dhcp lease for a particular interface
func Obtain(ifname string) error {
	cf, pf, lf := dhclientPathnames(ifname)
	pidData, err := ioutil.ReadFile(pf)
	if err != nil {
		return err
	}
	pidString := string(pidData)
	pid, err := strconv.Atoi(pidString)
	if err != nil {
		return err
	}
	unix.Kill(pid, unix.SIGTERM)
	os.Remove(pf)
	cmd := exec.Command("/sbin/dhclient", "-q", "-nw", "-cf", cf, "-pf", pf, "-lf", lf, ifname)
	return cmd.Run()
}

// Release releases a dhcp lease for a particular interface
func Release(ifname string) error {
	cf, pf, lf := dhclientPathnames(ifname)
	os.Remove(pf)
	cmd := exec.Command("/sbin/dhclient", "-q", "-cf", cf, "-pf", pf, "-lf", lf, "-r", ifname)
	return cmd.Run()
}

// Restart releases and then obtains a dhcp lease  for a particular interface
func Restart(ifname string) error {
	err := Release(ifname)
	if err != nil {
		return err
	}
	return Obtain(ifname)
}
