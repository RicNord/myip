package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Alias struct {
	Alias string `mapstructure:"alias"`
	IP    string `mapstructure:"ip"`
}

var aliasFlag bool

func getIp() (string, error) {
	resp, err := http.Get("https://icanhazip.com")
	if err != nil {
		return "", fmt.Errorf("failed to fetch IP address: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	ip := strings.TrimSuffix(string(body), "\n")
	return ip, nil
}

func getAliasOrIp(ip string) (ipOrAlias string, err error) {
	var aliases []Alias
	if err := viper.UnmarshalKey("aliases", &aliases); err != nil {
		return "", fmt.Errorf("unable to decode config into struct: %w", err)
	}

	for _, alias := range aliases {
		if alias.IP == ip {
			return alias.Alias, nil
		}
	}
	return ip, nil
}

func writeIpToFile(ip string) (err error) {
	tmpDir := os.TempDir()
	filePath := filepath.Join(tmpDir, ".my-last-known-ip")
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.WriteString(ip)
	if err != nil {
		return err
	}
	return nil
}

func getLastKownIp() (ip string, ok bool) {
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
		ip, ok := getLastKownIp()
		if ok {
			fmt.Printf("%s\n", ip)
		} else {
			fmt.Printf("No stored ip or alias exist\n")
		}
	},
}

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get your public IP address",
	Run: func(cmd *cobra.Command, args []string) {
		ip, err := getIp()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if aliasFlag {
			aliasOrIp, err := getAliasOrIp(ip)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			err = writeIpToFile(aliasOrIp)
			if err != nil {
				log.Println("Error:", err)
			}
			fmt.Printf("%s\n", aliasOrIp)
		} else {
			err := writeIpToFile(ip)
			if err != nil {
				log.Println("Error:", err)
			}
			fmt.Printf("%s\n", ip)
		}
	},
}

func init() {
	ipCmd.Flags().BoolVarP(&aliasFlag, "alias", "a", false, "Show the alias for the public IP address")
}
