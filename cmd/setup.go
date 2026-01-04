package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caravelcommerce/deck/internal/config"
	"github.com/caravelcommerce/deck/internal/docker"
	"github.com/caravelcommerce/deck/internal/traefik"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup Docker environment for the Magento project",
	Long:  `Reads deck.yaml and generates all Docker configuration files in the .deck folder.`,
	RunE:  runSetup,
}

func runSetup(cmd *cobra.Command, args []string) error {
	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	configPath := filepath.Join(cwd, "deck.yaml")
	deckDir := filepath.Join(cwd, ".deck")

	// Verifica se o .deck já existe
	if _, err := os.Stat(deckDir); err == nil {
		fmt.Println("⚠️  The .deck directory already exists.")
		if !askConfirmation("Do you want to overwrite it?") {
			fmt.Println("Setup cancelled.")
			return nil
		}
		fmt.Println("Removing existing .deck directory...")
		if err := os.RemoveAll(deckDir); err != nil {
			return fmt.Errorf("failed to remove .deck directory: %w", err)
		}
	}

	// Verifica se o deck.yaml existe
	var cfg *config.DeckConfig
	if !config.DeckYAMLExists(configPath) {
		fmt.Println("📋 deck.yaml not found. Attempting to detect Magento version from composer.json...")

		// Tenta detectar a versão do Magento
		magentoVersion, err := config.DetectMagentoVersion(cwd)
		if err != nil {
			return fmt.Errorf("failed to detect Magento version: %w\n\nPlease create a deck.yaml file manually with the Magento version", err)
		}

		fmt.Printf("✅ Detected Magento version: %s\n", magentoVersion)

		// Obtém o nome do projeto do diretório atual
		projectName := filepath.Base(cwd)

		fmt.Printf("📝 Creating deck.yaml with project name '%s' and Magento version '%s'...\n", projectName, magentoVersion)

		// Cria o deck.yaml
		if err := config.CreateDeckYAML(configPath, projectName, magentoVersion); err != nil {
			return fmt.Errorf("failed to create deck.yaml: %w", err)
		}

		fmt.Println("✅ deck.yaml created successfully!")
	}

	// Carrega o deck.yaml
	cfg, err = config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Printf("🚀 Setting up Deck environment for project: %s\n", cfg.Project)

	// Display configuration
	if cfg.Magento != "" {
		fmt.Printf("\n📦 Magento Version: %s\n", cfg.Magento)
		fmt.Println("   Using auto-detected versions:")
	} else {
		fmt.Println("\n📦 Using specified versions:")
	}
	fmt.Printf("   • PHP: %s\n", cfg.GetPHPVersion())
	if len(cfg.GetPHPExtensions()) > 0 {
		fmt.Printf("     Extensions: %v\n", cfg.GetPHPExtensions())
	}
	fmt.Printf("   • Nginx: %s\n", cfg.GetNginxVersion())
	fmt.Printf("   • MariaDB: %s\n", cfg.GetMariaDBVersion())
	fmt.Printf("   • OpenSearch: %s\n", cfg.GetOpenSearchVersion())
	fmt.Printf("   • Redis: %s\n", cfg.GetRedisVersion())
	fmt.Printf("   • RabbitMQ: %s\n", cfg.GetRabbitMQVersion())
	if cfg.IsNodeEnabled() {
		fmt.Printf("   • Node.js: %s\n", cfg.GetNodeVersion())
	}
	if cfg.IsSwooleEnabled() {
		fmt.Println("   • Swoole: enabled")
		if cfg.GetSwoolePort() > 0 {
			fmt.Printf("     API: https://api.%s.test (port %d)\n", cfg.Project, cfg.GetSwoolePort())
		}
	}
	fmt.Println()

	// Setup Traefik if not running
	if !traefik.IsTraefikRunning() {
		fmt.Println("📦 Setting up Traefik reverse proxy...")
		if err := traefik.SetupTraefik(); err != nil {
			return fmt.Errorf("failed to setup Traefik: %w", err)
		}
		fmt.Println("✅ Traefik is running")
	} else {
		fmt.Println("✅ Traefik is already running")
	}

	// Create .deck directory
	deckDir := filepath.Join(cwd, ".deck")
	if err := os.MkdirAll(deckDir, 0755); err != nil {
		return fmt.Errorf("failed to create .deck directory: %w", err)
	}

	// Generate Docker files
	fmt.Println("📝 Generating Docker configuration files...")
	if err := docker.GenerateDockerFiles(cfg, deckDir); err != nil {
		return fmt.Errorf("failed to generate Docker files: %w", err)
	}

	// Add .deck to .gitignore if it exists
	gitignorePath := filepath.Join(cwd, ".gitignore")
	if _, err := os.Stat(gitignorePath); err == nil {
		content, err := os.ReadFile(gitignorePath)
		if err == nil {
			gitignoreContent := string(content)
			if gitignoreContent == "" || gitignoreContent[len(gitignoreContent)-1] != '\n' {
				gitignoreContent += "\n"
			}
			if !contains(gitignoreContent, ".deck") {
				gitignoreContent += ".deck/\n"
				if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
					fmt.Printf("⚠️  Warning: failed to update .gitignore: %v\n", err)
				} else {
					fmt.Println("✅ Added .deck/ to .gitignore")
				}
			}
		}
	}

	fmt.Println("\n✨ Setup completed successfully!")
	fmt.Printf("\nYour project will be available at: https://%s.test\n", cfg.Project)
	if cfg.SwoolePort > 0 {
		fmt.Printf("Swoole API will be available at: https://api.%s.test\n", cfg.Project)
	}
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Run 'deck start' to start the environment")
	fmt.Printf("  2. Access your site at https://%s.test\n", cfg.Project)
	if cfg.SwoolePort > 0 {
		fmt.Printf("  3. Start Swoole server: deck bin/magento swoole:server:start\n")
		fmt.Printf("  4. Access Swoole API at https://api.%s.test\n", cfg.Project)
		fmt.Printf("  5. Run 'deck bin/magento' to execute other Magento commands\n")
	} else {
		fmt.Println("  3. Run 'deck bin/magento' to execute Magento commands")
	}
	fmt.Println("\nNote: You may need to add the SSL certificate to your trusted certificates.")

	traefikDir, err := traefik.GetTraefikDir()
	if err == nil {
		fmt.Printf("Certificate location: %s/certs/local-cert.pem\n", traefikDir)
	}

	return nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// askConfirmation solicita confirmação do usuário
func askConfirmation(message string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (y/N): ", message)

	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
