package cmd

import (
	"fmt"
	"os"

	"github.com/dvcrn/chainenv/backend"
	"github.com/dvcrn/chainenv/logger"
	"github.com/spf13/cobra"
)

var (
	backendType string
	opVault     string
	debug       bool
	log         *logger.Logger
	version     = "dev"
)

var rootCmd = &cobra.Command{
	Use:     "chainenv",
	Short:   "chainenv - A tool for managing environment variables securely",
	Long:    `chainenv allows you to securely store and retrieve environment variables using different secure backends like macOS Keychain or 1Password.`,
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log = logger.NewLogger(debug)
		log.Debug("Using backend: %s", backendType)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getBackendWithType(backendType string) (backend.Backend, error) {
	opts := backend.WithLogger(log)

	switch backendType {
	case "keychain":
		b, err := backend.NewKeychainBackend()
		if err != nil {
			return nil, fmt.Errorf("keychain backend unavailable: %w", err)
		}
		return b, nil
	case "1password":
		cfg, err := loadConfig()
		if err != nil {
			return nil, fmt.Errorf("error loading config: %w", err)
		}
		if err := ensureOpServiceAccountToken(cfg); err != nil {
			return nil, err
		}
		return backend.NewOnePasswordBackend(opVault, opts), nil
	default:
		return nil, fmt.Errorf("unknown backend: %s", backendType)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&backendType, "backend", "keychain", "Backend to use (keychain or 1password)")
	rootCmd.PersistentFlags().StringVar(&opVault, "vault", "chainenv", "1Password vault to use")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "Enable debug logging")
}
