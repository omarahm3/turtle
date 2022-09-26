package helpers

import (
	"fmt"
	"os/exec"
)

const sniffer = "nethogs"

var ErrVnstatDoesNotExist = fmt.Errorf("'%s' does not exist on this machine or not exposed on current user $PATH", sniffer)

func SnifferExists() (bool, error) {
	r := commandExists("nethogs")
	if !r {
		return false, ErrVnstatDoesNotExist
	}

	return true, nil
}

func commandExists(c string) bool {
	_, err := exec.LookPath(c)
	return err == nil
}
