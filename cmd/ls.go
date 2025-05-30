package cmd

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:     "ls",
	Aliases: []string{"list"},
	Short:   "List all stored accounts",
	Long:    `List all accounts that have passwords stored in the configured backend.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Debug("Listing all accounts")

		b, err := getBackendWithType(backendType)
		if err != nil {
			log.Err("Error initializing backend: %v", err)
			os.Exit(1)
		}

		accounts, err := b.List()
		if err != nil {
			log.Err("Error listing accounts: %v", err)
			os.Exit(1)
		}

		if len(accounts) == 0 {
			fmt.Println("No accounts found")
			return
		}

		// Sort accounts alphabetically
		sort.Strings(accounts)

		for _, account := range accounts {
			fmt.Println(account)
		}
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}