package cmd

import (
	"context"
	"github.com/mnizarzr/dot-test/db/seeder"
	"log"
	"time"

	"github.com/mnizarzr/dot-test/config"
	"github.com/mnizarzr/dot-test/db"
	"github.com/spf13/cobra"
)

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed the database with default data",
	Long:  `Seed the database with default users, projects, and other sample data for testing`,
	Run: func(cmd *cobra.Command, args []string) {
		seedDatabase()
	},
}

func seedDatabase() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database, err := db.NewPostgresGormDb(cfg.PgUri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	seederInstance := seeder.NewSeeder(database)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if err := seederInstance.SeedAll(ctx); err != nil {
		log.Fatalf("Failed to seed database: %v", err)
	}
}
