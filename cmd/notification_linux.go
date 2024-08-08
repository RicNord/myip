//go:build linux

package cmd

import (
	"os/exec"
)

func notificationService(ipOrAlias string) error {
	cmd := exec.Command("notify-send", "New IP detected:", ipOrAlias)
	if _, err := cmd.Output(); err != nil {
		return err
	}
	return nil
}
