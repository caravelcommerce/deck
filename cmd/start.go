package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caravelcommerce/deck/internal/config"
	"github.com/caravelcommerce/deck/internal/traefik"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the Docker environment",
	Long:  `Starts all Docker containers for the Magento project.`,
	RunE:  runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Check if .deck directory exists
	deckDir := filepath.Join(cwd, ".deck")
	if _, err := os.Stat(deckDir); os.IsNotExist(err) {
		return fmt.Errorf(".deck directory not found. Please run 'deck setup' first")
	}

	// Load deck.yaml to get project name
	configPath := filepath.Join(cwd, "deck.yaml")
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Ensure Traefik is running
	if !traefik.IsTraefikRunning() {
		fmt.Println("📦 Starting Traefik reverse proxy...")
		if err := traefik.SetupTraefik(); err != nil {
			return fmt.Errorf("failed to start Traefik: %w", err)
		}
	}

	fmt.Printf("🚀 Starting Docker environment for: %s\n", cfg.Project)

	// Run docker compose up
	dockerCmd := exec.Command("docker", "compose", "up", "-d")
	dockerCmd.Dir = deckDir
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	if err := dockerCmd.Run(); err != nil {
		return fmt.Errorf("failed to start Docker containers: %w", err)
	}

	fmt.Println("\n✅ Environment started successfully!")
	fmt.Printf("\n🌐 Your site is available at: https://%s.test\n", cfg.Project)
	if cfg.SwoolePort > 0 {
		fmt.Printf("🚀 Swoole API endpoint: https://api.%s.test (port %d)\n", cfg.Project, cfg.SwoolePort)
		fmt.Printf("   Start with: deck bin/magento swoole:server:start\n")
	}
	fmt.Println("\nServices:")
	fmt.Printf("  - Web: https://%s.test\n", cfg.Project)
	if cfg.SwoolePort > 0 {
		fmt.Printf("  - Swoole API: https://api.%s.test\n", cfg.Project)
	}
	fmt.Println("  - Traefik Dashboard: http://localhost:8080")
	fmt.Printf("  - Database: %s_mariadb:3306 (user: magento, password: magento)\n", cfg.Project)
	fmt.Printf("  - Redis: %s_redis:6379\n", cfg.Project)
	fmt.Printf("  - OpenSearch: %s_opensearch:9200\n", cfg.Project)
	fmt.Printf("  - RabbitMQ: %s_rabbitmq:15672 (user: guest, password: guest)\n", cfg.Project)

	return nil
}
