package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func getLastKnownIp() (ip string, exist bool) {
	var lastIP string
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, ".my-last-known-ip")

	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Failed to open file: %s, error: %v", filePath, err)
		return "", false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		lastIP = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Failed to read IP from file: %s, error: %v", filePath, err)
		return "", false
	}

	return lastIP, true
}

var lastIpCmd = &cobra.Command{
	Use:   "last",
	Short: "Get last known ip or alias",
	Run: func(cmd *cobra.Command, args []string) {
		ip, exist := getLastKnownIp()
		if exist {
			fmt.Printf("%s\n", ip)
		} else {
			fmt.Printf("No stored ip or alias exist\n")
		}
	},
}
