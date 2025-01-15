package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/dvcrn/chainenv/backend"
	"github.com/spf13/cobra"
)

var (
	shellType string
	fishFlag  bool
	bashFlag  bool
	zshFlag   bool
)

func formatShellExports(accountsPasswords map[string]string, shell string) string {
	if len(accountsPasswords) == 0 {
		return ""
	}

	var exports []string
	for account, password := range accountsPasswords {
		var format string
		switch shell {
		case "fish":
			format = "set -x %s '%s'"
		case "bash", "zsh", "sh":
			format = "export %s='%s'"
		default:
			format = "%s='%s'"
		}
		exports = append(exports, fmt.Sprintf(format, account, password))
	}
	return strings.Join(exports, "\n")
}

func getMultiplePasswords(accounts []string, b backend.Backend) map[string]string {
	results := make(map[string]string)
	for _, account := range accounts {
		if password, err := b.GetPassword(account); err == nil {
			results[account] = password
		}
	}
	return results
}

var getEnvCmd = &cobra.Command{
	Use:   "get-env [account1,account2,...]",
	Short: "Get passwords as environment variables",
	Long: `Retrieve passwords for multiple accounts and format them as environment variables.
Multiple accounts should be provided as a comma-separated list, e.g.:
  chainenv get-env AWS_KEY,AWS_SECRET --shell fish`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		accounts := strings.Split(args[0], ",")
		// Handle legacy shell flags
		switch {
		case fishFlag:
			shellType = "fish"
		case bashFlag:
			shellType = "bash"
		case zshFlag:
			shellType = "zsh"
		}

		log.Debug("Getting passwords for accounts: %s, shell=%s", strings.Join(accounts, ", "), shellType)

		b, err := getBackend()
		if err != nil {
			log.Err("Error initializing backend: %v", err)
			os.Exit(1)
		}

		passwords := getMultiplePasswords(accounts, b)
		output := formatShellExports(passwords, shellType)
		if output == "" {
			fmt.Fprintln(os.Stderr, "No passwords found")
			os.Exit(1)
		}
		fmt.Println(output)
	},
}

func init() {
	// New style
	getEnvCmd.Flags().StringVar(&shellType, "shell", "plain", "Shell format (fish, bash, zsh)")

	// Legacy style
	getEnvCmd.Flags().BoolVar(&fishFlag, "fish", false, "Use fish shell format (legacy)")
	getEnvCmd.Flags().BoolVar(&bashFlag, "bash", false, "Use bash shell format (legacy)")
	getEnvCmd.Flags().BoolVar(&zshFlag, "zsh", false, "Use zsh shell format (legacy)")

	rootCmd.AddCommand(getEnvCmd)
}
