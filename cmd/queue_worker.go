package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var queueCmd = &cobra.Command{
	Use:   "start-queue",
	Short: "Manage the message queue",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("(placeholder)")
	},
}
