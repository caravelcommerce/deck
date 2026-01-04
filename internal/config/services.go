package config

// ServiceConfig representa configuração de um serviço com versão e configurações customizadas
type ServiceConfig struct {
	Version       string                 `yaml:"version"`
	Configuration map[string]interface{} `yaml:"configuration,omitempty"`
}

// PHPConfig configuração específica do PHP
type PHPConfig struct {
	Version    string   `yaml:"version"`
	Extensions []string `yaml:"extensions,omitempty"`
}

// NginxConfig configuração específica do Nginx
type NginxConfig struct {
	Version       string                 `yaml:"version"`
	Configuration map[string]interface{} `yaml:"configuration,omitempty"`
}

// MariaDBConfig configuração específica do MariaDB
type MariaDBConfig struct {
	Version       string                 `yaml:"version"`
	Configuration map[string]interface{} `yaml:"configuration,omitempty"`
}

// OpenSearchConfig configuração específica do OpenSearch
type OpenSearchConfig struct {
	Version       string                 `yaml:"version"`
	Configuration map[string]interface{} `yaml:"configuration,omitempty"`
}

// RedisConfig configuração específica do Redis
type RedisConfig struct {
	Version       string                 `yaml:"version"`
	Configuration map[string]interface{} `yaml:"configuration,omitempty"`
}

// RabbitMQConfig configuração específica do RabbitMQ
type RabbitMQConfig struct {
	Version       string                 `yaml:"version"`
	Configuration map[string]interface{} `yaml:"configuration,omitempty"`
}

// NodeConfig configuração específica do Node.js
type NodeConfig struct {
	Version string `yaml:"version"`
}

// SwooleConfig configuração específica do Swoole
type SwooleConfig struct {
	Enabled bool `yaml:"enabled"`
	Port    int  `yaml:"port,omitempty"`
}

// Helper functions para obter versões
func (p *PHPConfig) GetVersion() string {
	if p != nil && p.Version != "" {
		return p.Version
	}
	return ""
}

func (n *NginxConfig) GetVersion() string {
	if n != nil && n.Version != "" {
		return n.Version
	}
	return ""
}

func (m *MariaDBConfig) GetVersion() string {
	if m != nil && m.Version != "" {
		return m.Version
	}
	return ""
}

func (o *OpenSearchConfig) GetVersion() string {
	if o != nil && o.Version != "" {
		return o.Version
	}
	return ""
}

func (r *RedisConfig) GetVersion() string {
	if r != nil && r.Version != "" {
		return r.Version
	}
	return ""
}

func (r *RabbitMQConfig) GetVersion() string {
	if r != nil && r.Version != "" {
		return r.Version
	}
	return ""
}

func (n *NodeConfig) GetVersion() string {
	if n != nil && n.Version != "" {
		return n.Version
	}
	return ""
}

// HasExtension verifica se uma extensão PHP está habilitada
func (p *PHPConfig) HasExtension(ext string) bool {
	if p == nil || p.Extensions == nil {
		return false
	}
	for _, e := range p.Extensions {
		if e == ext {
			return true
		}
	}
	return false
}

// GetConfigValue retorna um valor de configuração
func (s *ServiceConfig) GetConfigValue(key string) interface{} {
	if s == nil || s.Configuration == nil {
		return nil
	}
	return s.Configuration[key]
}

func (n *NginxConfig) GetConfigValue(key string) interface{} {
	if n == nil || n.Configuration == nil {
		return nil
	}
	return n.Configuration[key]
}

func (m *MariaDBConfig) GetConfigValue(key string) interface{} {
	if m == nil || m.Configuration == nil {
		return nil
	}
	return m.Configuration[key]
}

func (o *OpenSearchConfig) GetConfigValue(key string) interface{} {
	if o == nil || o.Configuration == nil {
		return nil
	}
	return o.Configuration[key]
}

func (r *RedisConfig) GetConfigValue(key string) interface{} {
	if r == nil || r.Configuration == nil {
		return nil
	}
	return r.Configuration[key]
}

func (r *RabbitMQConfig) GetConfigValue(key string) interface{} {
	if r == nil || r.Configuration == nil {
		return nil
	}
	return r.Configuration[key]
}
