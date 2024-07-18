package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
)

var monitorIpCmd = &cobra.Command{
	Use:   "monitor-ip",
	Short: "Monitor for IP changes and notify when a new IP or alias is detected",
	Run: func(cmd *cobra.Command, args []string) {
		for {
			currentIp, err := getIp()
			if err != nil {
				log.Printf("Error fetching current IP: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			currentIpOrAlias, err := getAliasOrIp(currentIp)
			if err != nil {
				log.Printf("Error getting alias or Ip for current value: %v", err)
			}

			lastIP, exists := getLastKnownIp()
			if !exists || currentIpOrAlias != lastIP {

				err = writeIpToFile(currentIpOrAlias)
				if err != nil {
					log.Printf("Error saving last known IP: %v", err)
				}

				err = notificationService(currentIpOrAlias)
				if err != nil {
					log.Printf("Error notifying new IP: %v", err)
				}
			}

			time.Sleep(5 * time.Second)
		}
	},
}
