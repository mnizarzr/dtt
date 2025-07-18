package cmd

import (
	"github.com/mnizarzr/dot-test/app"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	Run: func(cmd *cobra.Command, args []string) {
		app.Setup()
	},
}
