package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type Alias struct {
	Alias string `mapstructure:"alias"`
	IP    string `mapstructure:"ip"`
}

var cfgFile string
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

var rootCmd = &cobra.Command{
	Use:   "myip",
	Short: "A CLI tool to get public IP and alias for known IPs",
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.myip.json)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolVarP(&aliasFlag, "alias", "a", false, "Show the alias for the public IP address")

	rootCmd.AddCommand(lastIpCmd)

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".myip" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("json")
		viper.SetConfigName(".myip")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		// fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
