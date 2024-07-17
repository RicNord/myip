package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
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

func getAliasOrIp(ip string) (string, error) {
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
			fmt.Printf("%s\n", aliasOrIp)
		} else {
			fmt.Printf("%s\n", ip)
		}
	},
}

func init() {
	ipCmd.Flags().BoolVarP(&aliasFlag, "alias", "a", false, "Show the alias for the public IP address")
}
