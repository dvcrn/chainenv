package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dvcrn/chainenv/backend"
	"github.com/spf13/cobra"
)

var (
	generateEnvShellType string
	generateEnvFishFlag  bool
	generateEnvBashFlag  bool
	generateEnvZshFlag   bool
)

var generateEnvCmd = &cobra.Command{
	Use:   "generate-env",
	Short: "Generate environment exports from config",
	Long:  "Generate environment exports for keys declared in the chainenv config.",
	Run: func(cmd *cobra.Command, args []string) {
		// Handle legacy shell flags
		switch {
		case generateEnvFishFlag:
			generateEnvShellType = "fish"
		case generateEnvBashFlag:
			generateEnvShellType = "bash"
		case generateEnvZshFlag:
			generateEnvShellType = "zsh"
		}

		cfg, err := loadConfig()
		if err != nil {
			log.Err("Error loading config: %v", err)
			os.Exit(1)
		}
		if cfg == nil {
			fmt.Fprintln(os.Stderr, "No config found")
			os.Exit(1)
		}
		if len(cfg.Keys) == 0 {
			fmt.Fprintln(os.Stderr, "No keys found")
			return
		}

		var accounts []string
		for _, entry := range cfg.Keys {
			if entry.Name == "" {
				continue
			}
			accounts = append(accounts, entry.Name)
		}

		if len(accounts) == 0 {
			fmt.Fprintln(os.Stderr, "No keys found")
			return
		}

		backends := make(map[string]backend.Backend)
		getBackend := func(provider string) (backend.Backend, error) {
			if cached, ok := backends[provider]; ok {
				return cached, nil
			}
			b, err := getBackendWithType(provider)
			if err != nil {
				return nil, err
			}
			backends[provider] = b
			return b, nil
		}

		passwords := make(map[string]string)
		var firstErr error
		for _, account := range accounts {
			provider, defaultValue := resolveKeyConfig(cfg, account, backendType)
			b, err := getBackend(provider)
			if err != nil {
				if firstErr == nil {
					firstErr = err
				}
				continue
			}

			password, err := b.GetPassword(account)
			if err != nil {
				if errors.Is(err, backend.ErrNotFound) && defaultValue != nil {
					passwords[account] = *defaultValue
					continue
				}
				if firstErr == nil {
					firstErr = err
				}
				continue
			}
			passwords[account] = password
		}

		output := formatShellExports(passwords, generateEnvShellType)
		if output == "" {
			fmt.Fprintln(os.Stderr, "No passwords found")
			if firstErr != nil {
				fmt.Fprintln(os.Stderr, firstErr.Error())
			}
			os.Exit(1)
		}
		fmt.Println(output)
	},
}

func init() {
	// New style
	generateEnvCmd.Flags().StringVar(&generateEnvShellType, "shell", "plain", "Shell format (fish, bash, zsh)")

	// Legacy style
	generateEnvCmd.Flags().BoolVar(&generateEnvFishFlag, "fish", false, "Use fish shell format (legacy)")
	generateEnvCmd.Flags().BoolVar(&generateEnvBashFlag, "bash", false, "Use bash shell format (legacy)")
	generateEnvCmd.Flags().BoolVar(&generateEnvZshFlag, "zsh", false, "Use zsh shell format (legacy)")

	rootCmd.AddCommand(generateEnvCmd)
}
