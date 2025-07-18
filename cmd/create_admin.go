package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var createAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Create an admin user",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Admin user created (placeholder)")
	},
}
