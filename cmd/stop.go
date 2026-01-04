package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caravelcommerce/deck/internal/config"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the Docker environment",
	Long:  `Stops all Docker containers for the Magento project.`,
	RunE:  runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("🛑 Stopping Docker environment for: %s\n", cfg.Project)

	// Run docker compose down
	dockerCmd := exec.Command("docker", "compose", "down")
	dockerCmd.Dir = deckDir
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr

	if err := dockerCmd.Run(); err != nil {
		return fmt.Errorf("failed to stop Docker containers: %w", err)
	}

	fmt.Println("✅ Environment stopped successfully!")

	return nil
}
