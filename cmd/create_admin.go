package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/mnizarzr/dot-test/db/seeder"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mnizarzr/dot-test/config"
	"github.com/mnizarzr/dot-test/db"
	"github.com/spf13/cobra"
)

var createAdminCmd = &cobra.Command{
	Use:   "create-admin",
	Short: "Create an admin user",
	Long:  `Create an admin user interactively by providing name, email, and password`,
	Run: func(cmd *cobra.Command, args []string) {
		createAdminUser()
	},
}

func createAdminUser() {
	fmt.Println("🔧 Creating Admin User")
	fmt.Println("======================")

	// Load configuration
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	database, err := db.NewPostgresGormDb(cfg.PgUri)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create seeder
	seederInstance := seeder.NewSeeder(database)

	// Get admin details interactively
	reader := bufio.NewReader(os.Stdin)

	// Get name
	fmt.Print("Enter admin name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read name: %v", err)
	}
	name = strings.TrimSpace(name)

	// Get email
	fmt.Print("Enter admin email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read email: %v", err)
	}
	email = strings.TrimSpace(email)

	// Get password (hidden input)
	fmt.Print("Enter admin password: ")
	password, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read password: %v", err)
	}
	password = strings.TrimSpace(password)

	// Validate inputs
	if name == "" {
		log.Fatal("Name cannot be empty")
	}
	if email == "" {
		log.Fatal("Email cannot be empty")
	}
	if password == "" {
		log.Fatal("Password cannot be empty")
	}

	// Create admin user
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	adminUser, err := seederInstance.CreateAdminUser(ctx, name, email, password)
	if err != nil {
		log.Fatalf("Failed to create admin user: %v", err)
	}

	fmt.Printf("\n✅ Admin user created successfully!\n")
	fmt.Printf("   ID: %s\n", adminUser.ID)
	fmt.Printf("   Name: %s\n", adminUser.Name)
	fmt.Printf("   Email: %s\n", adminUser.Email)
	fmt.Printf("   Role: %s\n", adminUser.Role)
	fmt.Printf("   Created At: %s\n", adminUser.CreatedAt.Format("2006-01-02 15:04:05"))
}
