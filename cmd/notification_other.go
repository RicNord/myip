//go:build !linux

package cmd

import (
	"fmt"
)

func notificationService(ipOrAlias string) error {
	fmt.Println(ipOrAlias)
	return nil
}
