package cmd

import (
	"fmt"
	"os"

	"github.com/dvcrn/chainenv/config"
	"github.com/spf13/cobra"
)

var setDefault string

var setCmd = &cobra.Command{
	Use:   "set [account] [password]",
	Short: "Set a password for an account",
	Long:  `Store a new password for the specified account.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		account := args[0]
		password := args[1]
		log.Debug("Setting password for account: %s", account)

		b, err := getBackendWithType(backendType)
		if err != nil {
			log.Err("Error initializing backend: %v", err)
			os.Exit(1)
		}

		if err := b.SetPassword(account, password, false); err != nil {
			log.Err("Failed to set password: %v", err)
			os.Exit(1)
		}

		cwd, err := os.Getwd()
		if err != nil {
			log.Err("Failed to determine current directory: %v", err)
			os.Exit(1)
		}

		configPath, ok, err := config.FindConfig(cwd)
		if err != nil {
			log.Err("Failed to locate config file: %v", err)
			os.Exit(1)
		}
		if !ok {
			configPath = config.DefaultConfigPath(cwd)
		}

		cfg, err := config.LoadOrEmpty(configPath)
		if err != nil {
			log.Err("Failed to read config: %v", err)
			os.Exit(1)
		}

		entry := config.KeyEntry{Name: account, Provider: backendType}
		if existing, ok := cfg.FindKey(account); ok {
			entry.Default = existing.Default
		}
		if cmd.Flags().Changed("default") {
			entry.Default = &setDefault
		}
		cfg.UpsertKey(entry)

		if err := config.Save(configPath, cfg); err != nil {
			log.Err("Failed to write config: %v", err)
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

		b, err := getBackendWithType(backendType)
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
	setCmd.Flags().StringVar(&setDefault, "default", "", "Default value to store in config if secret is missing")
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(updateCmd)
}
