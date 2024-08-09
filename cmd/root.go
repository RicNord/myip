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
	url := viper.GetString("url")
	resp, err := http.Get(url)
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

	ip := strings.TrimSpace(string(body))
	return ip, nil
}

func getAliasOrIp(ip string) (ipOrAlias string, err error) {
	aliases, err := getAliases()
	if err != nil {
		return "", err
	}

	for _, alias := range aliases {
		if alias.IP == ip {
			return alias.Alias, nil
		}
	}
	return ip, nil
}

func writeIpToFile(ip string) (err error) {
	filePath := viper.GetString("state-file")
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %s, error: %v", filePath, err)
	}

	defer f.Close()
	_, err = f.WriteString(ip)
	if err != nil {
		return fmt.Errorf("failed to write IP to file: %s, error: %v", filePath, err)
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

		ipOrAlias := ip
		if aliasFlag {
			ipOrAlias, err = getAliasOrIp(ip)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
		}

		if err := writeIpToFile(ipOrAlias); err != nil {
			log.Printf("Error: %v", err)
		}
		fmt.Println(ipOrAlias)
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
	rootCmd.AddCommand(listAliasesCmd)
	rootCmd.AddCommand(monitorIpCmd)

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
	viper.SetDefault("url", "https://icanhazip.com")
	tempDir := os.TempDir()
	viper.SetDefault("state-file", filepath.Join(tempDir, ".my-last-known-ip"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		//log.Printf("Using config file: %s", viper.ConfigFileUsed())
	}
}
