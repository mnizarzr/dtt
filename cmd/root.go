package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dtt (dot test task)",
	Short: "run http or other command",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(queueCmd)
	rootCmd.AddCommand(createAdminCmd)
	rootCmd.AddCommand(seedCmd)
}
