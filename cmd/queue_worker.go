package cmd

import (
	"log"

	"github.com/mnizarzr/dot-test/config"
	"github.com/mnizarzr/dot-test/jobs"
	"github.com/spf13/cobra"
)

var queueCmd = &cobra.Command{
	Use:   "start-queue",
	Short: "Start the job queue worker",
	Run: func(cmd *cobra.Command, args []string) {

		config, err := config.LoadConfig(".")
		if err != nil {
			log.Fatal("Error loading config:", err)
		}

		jobManager := jobs.NewJobManager(config)

		log.Println("Starting job queue worker...")
		if err := jobManager.Start(); err != nil {
			log.Fatal("Error starting job manager:", err)
		}

	},
}
