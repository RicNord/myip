package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"strings"
)

var aliasFlag bool

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "Get your public IP address",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get("https://icanhazip.com")
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Printf("error: status code: %v\n", resp.StatusCode)
			return
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error: %v\n", err)
			return
		}

		ip := strings.TrimSuffix(string(body), "\n")

		if aliasFlag {
			aliases := viper.Get("aliases").([]interface{})
			for _, a := range aliases {
				alias := a.(map[string]interface{})
				if alias["ip"] == ip {
					fmt.Printf("%s\n", alias["alias"])
					return
				}
			}
			fmt.Printf("%s\n", ip)
		} else {
			fmt.Printf("%s\n", ip)
		}
	},
}

func init() {
	ipCmd.Flags().BoolVarP(&aliasFlag, "alias", "a", false, "Show the alias for the public IP address")
}
