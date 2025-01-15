package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dvcrn/chainenv/backend"
	"github.com/dvcrn/chainenv/logger"
)

var (
	backendType = flag.String("backend", "keychain", "Backend to use (keychain or 1password)")
	opVault     = flag.String("vault", "chainenv", "1Password vault to use")
)

func getBackend(opts ...backend.BackendOption) (backend.Backend, error) {
	switch *backendType {
	case "keychain":
		return backend.NewKeychainBackend(), nil
	case "1password":
		return backend.NewOnePasswordBackend(*opVault, opts...), nil
	default:
		return nil, fmt.Errorf("unknown backend: %s", *backendType)
	}
}

func getMultiplePasswords(b backend.Backend, accounts []string) map[string]string {
	results := make(map[string]string)
	for _, account := range accounts {
		if password, err := b.GetPassword(account); err == nil {
			results[account] = password
		}
	}
	return results
}

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

func main() {
	debug := true
	flag.Parse()
	args := flag.Args()

	logger := logger.NewLogger(debug)

	if len(args) < 1 {
		logger.Err("Expected 'get', 'get-env', 'set', or 'update' subcommands")
		os.Exit(1)
	}

	logger.Debug("Using backend: %s", *backendType)

	backendOpts := backend.WithLogger(logger)
	b, err := getBackend(backendOpts)
	if err != nil {
		logger.Err("Error initializing backend: %v\n", err)
		os.Exit(1)
	}

	switch args[0] {
	case "get":
		logger.Debug("Running 'get'")

		var account string
		if len(args) >= 2 && !strings.HasPrefix(args[1], "-") {
			// Positional argument style
			account = args[1]
		} else {
			// Flag style
			getCmd := flag.NewFlagSet("get", flag.ExitOnError)
			getAccount := getCmd.String("account", "", "Account name to retrieve")
			getCmd.Parse(args[1:])
			account = *getAccount
		}

		logger.Debug("Getting password for account: %s", account)

		if account == "" {
			logger.Err("account is required")
			os.Exit(1)
		}

		if password, err := b.GetPassword(account); err != nil {
			logger.Err("Error retrieving password: %v", err)
			os.Exit(1)
		} else {
			fmt.Println(password)
		}

	case "get-env":
		logger.Debug("Running 'get-env'")

		var accounts []string
		var shell string = "plain"

		if len(args) >= 2 && !strings.HasPrefix(args[1], "-") {
			// Positional argument style
			accounts = strings.Split(args[1], ",")
			// Check for shell flags
			for i := 2; i < len(args); i++ {
				switch args[i] {
				case "--fish":
					shell = "fish"
				case "--bash":
					shell = "bash"
				case "--zsh":
					shell = "zsh"
				}
			}

		} else {
			// Flag style
			getEnvCmd := flag.NewFlagSet("get-env", flag.ExitOnError)
			getEnvAccounts := getEnvCmd.String("accounts", "", "Comma-separated list of account names to retrieve")
			fishFlag := getEnvCmd.Bool("fish", false, "Output fish shell format")
			bashFlag := getEnvCmd.Bool("bash", false, "Output bash shell format")
			zshFlag := getEnvCmd.Bool("zsh", false, "Output zsh shell format")
			getEnvCmd.Parse(args[1:])

			if *getEnvAccounts != "" {
				accounts = strings.Split(*getEnvAccounts, ",")
			}

			switch {
			case *fishFlag:
				shell = "fish"
			case *bashFlag:
				shell = "bash"
			case *zshFlag:
				shell = "zsh"
			}
		}

		if len(accounts) == 0 {
			logger.Err("accounts are required")
			os.Exit(1)
		}

		logger.Debug("Getting passwords for accounts: %s, shell=%s", strings.Join(accounts, ", "), shell)

		passwords := getMultiplePasswords(b, accounts)
		output := formatShellExports(passwords, shell)
		if output == "" {
			fmt.Fprintln(os.Stderr, "No passwords found")
			os.Exit(1)
		}
		fmt.Println(output)

	case "set", "update":
		var account, password string
		isUpdate := args[0] == "update"

		if isUpdate {
			logger.Debug("Running 'update'")
		} else {
			logger.Debug("Running 'set'")
		}

		if len(args) >= 3 && !strings.HasPrefix(args[1], "-") {
			// Positional argument style
			account = args[1]
			password = args[2]
		} else {
			// Flag style
			cmd := flag.NewFlagSet(args[0], flag.ExitOnError)
			accountFlag := cmd.String("account", "", "Account name")
			passFlag := cmd.String("password", "", "Password to store")
			cmd.Parse(args[1:])
			account = *accountFlag
			password = *passFlag
		}

		if account == "" || password == "" {
			fmt.Println("Error: both account and password are required")
			os.Exit(1)
		}

		logger.Debug("Setting password for account: %s to %s", account, password)

		if err := b.SetPassword(account, password, isUpdate); err != nil {
			logger.Err("Failed to %s password: %v\n", args[0], err)
			os.Exit(1)
		}
		fmt.Printf("Password %s for %s\n", map[bool]string{true: "updated", false: "set"}[isUpdate], account)

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		fmt.Println("Expected 'get', 'get-env', 'set', or 'update' subcommands")
		os.Exit(1)
	}
}
