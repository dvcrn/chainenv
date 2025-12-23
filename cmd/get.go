package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dvcrn/chainenv/backend"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get [account]",
	Short: "Get a password for an account",
	Long:  `Retrieve a password stored for the specified account.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		log.Debug("Getting password for account: %s", account)

		cfg, err := loadConfig()
		if err != nil {
			log.Err("Error loading config: %v", err)
			os.Exit(1)
		}

		provider, defaultValue := resolveKeyConfig(cfg, account, backendType)
		b, err := getBackendWithType(provider)
		if err != nil {
			log.Err("Error initializing backend: %v", err)
			os.Exit(1)
		}

		password, err := b.GetPassword(account)
		if err != nil {
			if errors.Is(err, backend.ErrNotFound) && defaultValue != nil {
				fmt.Println(*defaultValue)
				return
			}
			log.Err("Error retrieving password: %v", err)
			os.Exit(1)
		}

		fmt.Println(password)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
