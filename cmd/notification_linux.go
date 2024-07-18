//go:build linux

package cmd

import (
	"os/exec"
)

func notificationService(ipOrAlias string) error {
	app := "notify-send"

	arg0 := "New IP detected:"
	arg1 := ipOrAlias

	cmd := exec.Command(app, arg0, arg1)
	_, err := cmd.Output()

	if err != nil {
		return err
	}
	return nil
}
