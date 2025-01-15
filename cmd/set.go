package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set [account] [password]",
	Short: "Set a password for an account",
	Long:  `Store a new password for the specified account.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		password := args[1]
		log.Debug("Setting password for account: %s", account)

		b, err := getBackend()
		if err != nil {
			log.Err("Error initializing backend: %v", err)
			os.Exit(1)
		}

		if err := b.SetPassword(account, password, false); err != nil {
			log.Err("Failed to set password: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Password set for %s\n", account)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [account] [password]",
	Short: "Update a password for an existing account",
	Long:  `Update the password for an existing account.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		password := args[1]
		log.Debug("Updating password for account: %s", account)

		b, err := getBackend()
		if err != nil {
			log.Err("Error initializing backend: %v", err)
			os.Exit(1)
		}

		if err := b.SetPassword(account, password, true); err != nil {
			log.Err("Failed to update password: %v", err)
			os.Exit(1)
		}

		fmt.Printf("Password updated for %s\n", account)
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(updateCmd)
}
