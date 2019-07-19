package vyatta

import (
	"os/exec"
)

// Run runs vyatta-interfaces
func Run(ifname string, args ...string) error {
	cmdArgs := append([]string{"--dev", ifname}, args...)
	cmd := exec.Command("/opt/vyatta/sbin/vyatta-interfaces.pl", cmdArgs...)
	return cmd.Run()
}
