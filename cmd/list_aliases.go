package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func getAliases() ([]Alias, error) {
	var aliases []Alias
	if err := viper.UnmarshalKey("aliases", &aliases); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}
	return aliases, nil
}

var listAliasesCmd = &cobra.Command{
	Use:   "list-aliases",
	Short: "List all known IP aliases",
	Run: func(cmd *cobra.Command, args []string) {
		aliases, err := getAliases()
		if err != nil {
			log.Printf("Error fetching aliases: %v", err)
			return
		}

		if len(aliases) == 0 {
			log.Println("No aliases found.")
			return
		}

		for _, alias := range aliases {
			fmt.Printf("IP: %s, Alias: %s\n", alias.IP, alias.Alias)
		}
	},
}
