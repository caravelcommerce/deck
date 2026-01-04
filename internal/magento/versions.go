package magento

import (
	"embed"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

//go:embed versions/*.yaml
var versionsFS embed.FS

// ServiceVersion representa uma versão de serviço no arquivo YAML
type ServiceVersion struct {
	Version string `yaml:"version"`
}

// MagentoVersion estrutura completa de uma versão do Magento
type MagentoVersion struct {
	Version    string          `yaml:"version"`
	PHP        *ServiceVersion `yaml:"php"`
	Nginx      *ServiceVersion `yaml:"nginx"`
	MariaDB    *ServiceVersion `yaml:"mariadb"`
	OpenSearch *ServiceVersion `yaml:"opensearch"`
	Redis      *ServiceVersion `yaml:"redis"`
	RabbitMQ   *ServiceVersion `yaml:"rabbitmq"`
}

// MagentoRequirements versão simplificada para backward compatibility
type MagentoRequirements struct {
	Version    string
	PHP        string
	Nginx      string
	MariaDB    string
	OpenSearch string
	Redis      string
	RabbitMQ   string
}

// versionCache armazena as versões carregadas em memória
var versionCache map[string]*MagentoVersion

// init carrega todas as versões na inicialização
func init() {
	versionCache = make(map[string]*MagentoVersion)
	loadAllVersions()
}

// loadAllVersions carrega todos os arquivos YAML do diretório versions
func loadAllVersions() {
	entries, err := versionsFS.ReadDir("versions")
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".yaml") {
			continue
		}

		data, err := versionsFS.ReadFile(filepath.Join("versions", entry.Name()))
		if err != nil {
			continue
		}

		var ver MagentoVersion
		if err := yaml.Unmarshal(data, &ver); err != nil {
			continue
		}

		// Armazena usando a versão como chave
		versionCache[ver.Version] = &ver
	}
}

// GetRequirements retorna os requisitos para uma versão do Magento
func GetRequirements(version string) (*MagentoRequirements, error) {
	ver := GetVersion(version)
	if ver == nil {
		return nil, fmt.Errorf("unsupported Magento version: %s (versions available: %s)",
			version, strings.Join(GetSupportedVersions(), ", "))
	}

	// Converte para MagentoRequirements (formato simplificado)
	req := &MagentoRequirements{
		Version: ver.Version,
	}

	if ver.PHP != nil {
		req.PHP = ver.PHP.Version
	}
	if ver.Nginx != nil {
		req.Nginx = ver.Nginx.Version
	}
	if ver.MariaDB != nil {
		req.MariaDB = ver.MariaDB.Version
	}
	if ver.OpenSearch != nil {
		req.OpenSearch = ver.OpenSearch.Version
	}
	if ver.Redis != nil {
		req.Redis = ver.Redis.Version
	}
	if ver.RabbitMQ != nil {
		req.RabbitMQ = ver.RabbitMQ.Version
	}

	return req, nil
}

// GetVersion retorna a versão completa do Magento
func GetVersion(version string) *MagentoVersion {
	// Tenta buscar a versão exata primeiro
	if ver, ok := versionCache[version]; ok {
		return ver
	}

	// Se não encontrar, tenta versão sem sufixo (ex: 2.4.8-p5 -> 2.4.8)
	baseVersion := getBaseVersion(version)
	if baseVersion != version {
		if ver, ok := versionCache[baseVersion]; ok {
			return ver
		}
	}

	return nil
}

// getBaseVersion extrai a versão base (ex: "2.4.8-p3" -> "2.4.8")
func getBaseVersion(version string) string {
	// Se começa com patch (-p), remove o patch
	if idx := strings.Index(version, "-p"); idx > 0 {
		return version[:idx]
	}
	return version
}

// GetSupportedVersions retorna uma lista ordenada de versões suportadas
func GetSupportedVersions() []string {
	versions := make([]string, 0, len(versionCache))
	for version := range versionCache {
		versions = append(versions, version)
	}
	return versions
}

// GetLatestVersion retorna a versão mais recente disponível
func GetLatestVersion() string {
	versions := GetSupportedVersions()
	if len(versions) == 0 {
		return ""
	}

	// Retorna a última versão (assumindo que os nomes dos arquivos estão ordenados)
	latest := versions[0]
	for _, v := range versions {
		if v > latest {
			latest = v
		}
	}
	return latest
}

// ListVersionsByMajorMinor agrupa versões por versão maior.menor
func ListVersionsByMajorMinor() map[string][]string {
	result := make(map[string][]string)

	for version := range versionCache {
		parts := strings.SplitN(version, ".", 3)
		if len(parts) >= 2 {
			majorMinor := parts[0] + "." + parts[1]
			result[majorMinor] = append(result[majorMinor], version)
		}
	}

	return result
}

