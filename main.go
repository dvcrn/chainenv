package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func getPassword(account string) (string, error) {
	cmd := exec.Command("security", "find-generic-password", "-a", account, "-s", fmt.Sprintf("chainenv-%s", account), "-w")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error retrieving password: %v", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func getMultiplePasswords(accounts []string) map[string]string {
	results := make(map[string]string)
	for _, account := range accounts {
		if password, err := getPassword(account); err == nil {
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

func setPassword(account, password string, update bool) error {
	args := []string{"add-generic-password", "-a", account, "-s", fmt.Sprintf("chainenv-%s", account), "-w", password, "-j", "Set by chainenv"}
	if update {
		args = append(args, "-U")
	}

	cmd := exec.Command("security", args...)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("error setting password: %v: %s", err, output)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected 'get', 'get-env', 'set', or 'update' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "get":
		var account string
		if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "-") {
			// Positional argument style
			account = os.Args[2]
		} else {
			// Flag style
			getCmd := flag.NewFlagSet("get", flag.ExitOnError)
			getAccount := getCmd.String("account", "", "Account name to retrieve")
			getCmd.Parse(os.Args[2:])
			account = *getAccount
		}

		if account == "" {
			fmt.Println("Error: account is required")
			os.Exit(1)
		}

		if password, err := getPassword(account); err != nil {
			fmt.Fprintf(os.Stderr, "Password not found: %v\n", err)
			os.Exit(1)
		} else {
			fmt.Println(password)
		}

	case "get-env":
		var accounts []string
		var shell string = "plain"

		if len(os.Args) >= 3 && !strings.HasPrefix(os.Args[2], "-") {
			// Positional argument style
			accounts = strings.Split(os.Args[2], ",")
			// Check for shell flags
			for i := 3; i < len(os.Args); i++ {
				switch os.Args[i] {
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
			getEnvCmd.Parse(os.Args[2:])

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
			fmt.Println("Error: accounts are required")
			os.Exit(1)
		}

		passwords := getMultiplePasswords(accounts)
		output := formatShellExports(passwords, shell)
		if output == "" {
			fmt.Fprintln(os.Stderr, "No passwords found")
			os.Exit(1)
		}
		fmt.Println(output)

	case "set", "update":
		var account, password string
		isUpdate := os.Args[1] == "update"

		if len(os.Args) >= 4 && !strings.HasPrefix(os.Args[2], "-") {
			// Positional argument style
			account = os.Args[2]
			password = os.Args[3]
		} else {
			// Flag style
			cmd := flag.NewFlagSet(os.Args[1], flag.ExitOnError)
			accountFlag := cmd.String("account", "", "Account name")
			passFlag := cmd.String("password", "", "Password to store")
			cmd.Parse(os.Args[2:])
			account = *accountFlag
			password = *passFlag
		}

		if account == "" || password == "" {
			fmt.Println("Error: both account and password are required")
			os.Exit(1)
		}

		if err := setPassword(account, password, isUpdate); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to %s password: %v\n", os.Args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Password %s for %s\n", map[bool]string{true: "updated", false: "set"}[isUpdate], account)

	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		fmt.Println("Expected 'get', 'get-env', 'set', or 'update' subcommands")
		os.Exit(1)
	}
}
