package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	source string
	target string
)

var copyCmd = &cobra.Command{
	Use:     "copy [keys...]",
	Aliases: []string{"cp"},
	Short:   "Copy passwords between backends",
	Long:    `Copy specified passwords from one backend to another (keychain <-> 1password)`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		sourceBackend, err := getBackendWithType(source)
		if err != nil {
			log.Err("Error initializing source backend: %v", err)
			os.Exit(1)
		}

		targetBackend, err := getBackendWithType(target)
		if err != nil {
			log.Err("Error initializing target backend: %v", err)
			os.Exit(1)
		}

		keys := []string{}
		for _, arg := range args {
			if strings.Contains(arg, ",") {
				keys = append(keys, strings.Split(arg, ",")...)
			} else {
				keys = append(keys, arg)
			}
		}

		for _, key := range keys {
			password, err := sourceBackend.GetPassword(key)
			if err != nil {
				log.Err("Failed to get password for %s: %v", key, err)
				continue
			}

			if err := targetBackend.SetPassword(key, password, false); err != nil {
				if err := targetBackend.SetPassword(key, password, true); err != nil {
					log.Err("Failed to copy password for %s: %v", key, err)
					continue
				}
			}

			fmt.Printf("Copied password for %s from %s to %s\n", key, source, target)
		}
	},
}

func init() {
	copyCmd.Flags().StringVar(&source, "from", "", "Source backend (keychain or 1password)")
	copyCmd.Flags().StringVar(&target, "to", "", "Target backend (keychain or 1password)")
	copyCmd.MarkFlagRequired("from")
	copyCmd.MarkFlagRequired("to")
	rootCmd.AddCommand(copyCmd)
}
