package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/caravelcommerce/deck/internal/config"
	"github.com/spf13/cobra"
)

var binMagentoCmd = &cobra.Command{
	Use:                "bin/magento",
	Short:              "Execute Magento CLI commands",
	Long:               `Runs bin/magento commands inside the PHP container.`,
	DisableFlagParsing: true,
	RunE:               runBinMagento,
}

func runBinMagento(cmd *cobra.Command, args []string) error {
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

	containerName := fmt.Sprintf("%s_php", cfg.Project)

	// Check if container is running
	checkCmd := exec.Command("docker", "ps", "--filter", fmt.Sprintf("name=%s", containerName), "--format", "{{.Names}}")
	output, err := checkCmd.Output()
	if err != nil || len(output) == 0 {
		return fmt.Errorf("PHP container is not running. Please run 'deck start' first")
	}

	// Build docker exec command
	dockerArgs := []string{"exec", "-it", containerName, "php", "bin/magento"}
	dockerArgs = append(dockerArgs, args...)

	// Execute bin/magento in the PHP container
	dockerCmd := exec.Command("docker", dockerArgs...)
	dockerCmd.Stdin = os.Stdin
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	dockerCmd.Dir = cwd

	if err := dockerCmd.Run(); err != nil {
		return fmt.Errorf("failed to execute bin/magento: %w", err)
	}

	return nil
}
