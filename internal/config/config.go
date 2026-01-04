package config

import (
	"fmt"
	"os"

	"github.com/caravelcommerce/deck/internal/magento"
	"gopkg.in/yaml.v3"
)

// DeckConfig estrutura principal de configuração
type DeckConfig struct {
	Project    string            `yaml:"project"`  // Nome do projeto
	Magento    string            `yaml:"magento"`  // Versão do Magento
	PHP        *PHPConfig        `yaml:"php,omitempty"`
	Nginx      *NginxConfig      `yaml:"nginx,omitempty"`
	MariaDB    *MariaDBConfig    `yaml:"mariadb,omitempty"`
	OpenSearch *OpenSearchConfig `yaml:"opensearch,omitempty"`
	Redis      *RedisConfig      `yaml:"redis,omitempty"`
	RabbitMQ   *RabbitMQConfig   `yaml:"rabbitmq,omitempty"`
	Node       *NodeConfig       `yaml:"node,omitempty"`
	Swoole     *SwooleConfig     `yaml:"swoole,omitempty"`
}

// LoadConfig carrega e processa a configuração
func LoadConfig(path string) (*DeckConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read deck.yaml: %w", err)
	}

	var config DeckConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse deck.yaml: %w", err)
	}

	// Validate required fields
	if config.Project == "" {
		return nil, fmt.Errorf("project name is required in deck.yaml")
	}

	// Apply Magento version defaults if specified
	if config.Magento != "" {
		if err := config.applyMagentoDefaults(); err != nil {
			return nil, err
		}
	}

	// Apply final defaults
	config.applyDefaults()

	return &config, nil
}

// applyMagentoDefaults aplica defaults baseados na versão do Magento
func (c *DeckConfig) applyMagentoDefaults() error {
	requirements, err := magento.GetRequirements(c.Magento)
	if err != nil {
		return fmt.Errorf("failed to get Magento requirements: %w", err)
	}

	// PHP
	if c.PHP == nil {
		c.PHP = &PHPConfig{}
	}
	if c.PHP.Version == "" {
		c.PHP.Version = requirements.PHP
	}

	// Nginx
	if c.Nginx == nil {
		c.Nginx = &NginxConfig{}
	}
	if c.Nginx.Version == "" {
		c.Nginx.Version = requirements.Nginx
	}

	// MariaDB
	if c.MariaDB == nil {
		c.MariaDB = &MariaDBConfig{}
	}
	if c.MariaDB.Version == "" {
		c.MariaDB.Version = requirements.MariaDB
	}

	// OpenSearch
	if c.OpenSearch == nil {
		c.OpenSearch = &OpenSearchConfig{}
	}
	if c.OpenSearch.Version == "" {
		c.OpenSearch.Version = requirements.OpenSearch
	}

	// Redis
	if c.Redis == nil {
		c.Redis = &RedisConfig{}
	}
	if c.Redis.Version == "" {
		c.Redis.Version = requirements.Redis
	}

	// RabbitMQ
	if c.RabbitMQ == nil {
		c.RabbitMQ = &RabbitMQConfig{}
	}
	if c.RabbitMQ.Version == "" {
		c.RabbitMQ.Version = requirements.RabbitMQ
	}

	return nil
}

// applyDefaults aplica defaults finais
func (c *DeckConfig) applyDefaults() {
	// PHP defaults
	if c.PHP == nil {
		c.PHP = &PHPConfig{Version: "8.3"}
	}
	if c.PHP.Version == "" {
		c.PHP.Version = "8.3"
	}
	// Extensões padrão do Magento
	if c.PHP.Extensions == nil || len(c.PHP.Extensions) == 0 {
		c.PHP.Extensions = []string{
			"bcmath", "gd", "intl", "mbstring", "pdo_mysql",
			"soap", "sockets", "xsl", "zip", "opcache",
		}
	}

	// Nginx defaults
	if c.Nginx == nil {
		c.Nginx = &NginxConfig{Version: "1.28"}
	}
	if c.Nginx.Version == "" {
		c.Nginx.Version = "1.28"
	}

	// MariaDB defaults
	if c.MariaDB == nil {
		c.MariaDB = &MariaDBConfig{Version: "11.4"}
	}
	if c.MariaDB.Version == "" {
		c.MariaDB.Version = "11.4"
	}

	// OpenSearch defaults
	if c.OpenSearch == nil {
		c.OpenSearch = &OpenSearchConfig{Version: "3"}
	}
	if c.OpenSearch.Version == "" {
		c.OpenSearch.Version = "3"
	}

	// Redis defaults
	if c.Redis == nil {
		c.Redis = &RedisConfig{Version: "7.4"}
	}
	if c.Redis.Version == "" {
		c.Redis.Version = "7.4"
	}

	// RabbitMQ defaults
	if c.RabbitMQ == nil {
		c.RabbitMQ = &RabbitMQConfig{Version: "4.1"}
	}
	if c.RabbitMQ.Version == "" {
		c.RabbitMQ.Version = "4.1"
	}

	// Swoole defaults
	if c.Swoole != nil && c.Swoole.Enabled && c.Swoole.Port == 0 {
		c.Swoole.Port = 9501
	}
}

// Helper methods
func (c *DeckConfig) GetPHPExtensions() []string {
	if c.PHP == nil || c.PHP.Extensions == nil {
		return []string{}
	}
	return c.PHP.Extensions
}

func (c *DeckConfig) HasPHPExtension(ext string) bool {
	if c.PHP == nil {
		return false
	}
	return c.PHP.HasExtension(ext)
}

func (c *DeckConfig) GetNodeVersion() string {
	return c.Node.GetVersion()
}

func (c *DeckConfig) IsNodeEnabled() bool {
	return c.Node != nil && c.Node.Version != ""
}

func (c *DeckConfig) IsSwooleEnabled() bool {
	return c.Swoole != nil && c.Swoole.Enabled
}

func (c *DeckConfig) GetSwoolePort() int {
	if c.Swoole == nil {
		return 0
	}
	return c.Swoole.Port
}

// GetPHPVersion retorna a versão do PHP
func (c *DeckConfig) GetPHPVersion() string {
	if c.PHP == nil {
		return ""
	}
	return c.PHP.Version
}

// GetNginxVersion retorna a versão do Nginx
func (c *DeckConfig) GetNginxVersion() string {
	if c.Nginx == nil {
		return ""
	}
	return c.Nginx.Version
}

// GetMariaDBVersion retorna a versão do MariaDB
func (c *DeckConfig) GetMariaDBVersion() string {
	if c.MariaDB == nil {
		return ""
	}
	return c.MariaDB.Version
}

// GetOpenSearchVersion retorna a versão do OpenSearch
func (c *DeckConfig) GetOpenSearchVersion() string {
	if c.OpenSearch == nil {
		return ""
	}
	return c.OpenSearch.Version
}

// GetRedisVersion retorna a versão do Redis
func (c *DeckConfig) GetRedisVersion() string {
	if c.Redis == nil {
		return ""
	}
	return c.Redis.Version
}

// GetRabbitMQVersion retorna a versão do RabbitMQ
func (c *DeckConfig) GetRabbitMQVersion() string {
	if c.RabbitMQ == nil {
		return ""
	}
	return c.RabbitMQ.Version
}

// CreateDeckYAML cria um arquivo deck.yaml com a configuração fornecida
func CreateDeckYAML(path, projectName, magentoVersion string) error {
	config := &DeckConfig{
		Project: projectName,
		Magento: magentoVersion,
	}

	// Serializa para YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Adiciona comentário no topo
	header := `# Deck - Magento 2 Development Environment
# Auto-generated configuration file

`
	content := header + string(data)

	// Escreve o arquivo
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write deck.yaml: %w", err)
	}

	return nil
}

// DeckYAMLExists verifica se o deck.yaml existe
func DeckYAMLExists(projectPath string) bool {
	_, err := os.Stat(projectPath)
	return err == nil
}
