package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/dvcrn/chainenv/backend"
	"github.com/spf13/cobra"
)

var diagCmd = &cobra.Command{
	Use:   "diag",
	Short: "Diagnose available backends",
	Long:  "Checks availability of supported backends on this system: macOS Keychain, Linux Secret Service keyring, and 1Password.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Backend diagnostics:")

		// macOS Keychain
		if runtime.GOOS == "darwin" {
			if _, err := exec.LookPath("security"); err != nil {
				fmt.Println("- macOS Keychain: unavailable (security CLI not found)")
			} else {
				// Try a lightweight command to ensure it's functional
				c := exec.Command("security", "list-keychains")
				if err := c.Run(); err != nil {
					fmt.Printf("- macOS Keychain: unavailable (%v)\n", err)
				} else {
					fmt.Println("- macOS Keychain: available")
				}
			}
		} else {
			fmt.Println("- macOS Keychain: unavailable (not macOS)")
		}

		// Linux Keychain (Secret Service/KWallet via keyring)
		if runtime.GOOS == "linux" {
			if _, err := backend.NewKeychainBackend(); err != nil {
				fmt.Printf("- Linux Keyring (Secret Service/KWallet): unavailable (%v)\n", err)
			} else {
				fmt.Println("- Linux Keyring (Secret Service/KWallet): available")
			}
		} else {
			fmt.Println("- Linux Keyring (Secret Service/KWallet): unavailable (not Linux)")
		}

		// 1Password CLI
		if _, err := exec.LookPath("op"); err != nil {
			fmt.Println("- 1Password CLI: unavailable (op CLI not found)")
		} else {
			// Attempt a lightweight identity check
			c := exec.Command("op", "whoami", "--format", "json")
			if err := c.Run(); err != nil {
				fmt.Println("- 1Password CLI: installed, but not signed in")
			} else {
				fmt.Println("- 1Password CLI: available (signed in)")
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(diagCmd)
}
