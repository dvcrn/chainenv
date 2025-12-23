package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys declared in config",
	Long:  "List keys declared in .chainenv.toml or chainenv.toml.",
	Run: func(cmd *cobra.Command, args []string) {
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
			fmt.Println("No keys found")
			return
		}

		for _, entry := range cfg.Keys {
			fmt.Println(entry.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
