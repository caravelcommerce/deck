package traefik

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const traefikDockerCompose = `version: '3.8'

services:
  traefik:
    image: traefik:v3.0
    container_name: deck_traefik
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--providers.file.directory=/etc/traefik/dynamic"
      - "--providers.file.watch=true"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--entrypoints.web.http.redirections.entrypoint.to=websecure"
      - "--entrypoints.web.http.redirections.entrypoint.scheme=https"
    ports:
      - "80:80"
      - "443:443"
      - "8080:8080"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./certs:/etc/traefik/certs:ro
      - ./dynamic:/etc/traefik/dynamic:ro
    networks:
      - traefik_network
    restart: unless-stopped

networks:
  traefik_network:
    name: traefik_network
    driver: bridge
`

const traefikDynamicConfig = `tls:
  certificates:
    - certFile: /etc/traefik/certs/local-cert.pem
      keyFile: /etc/traefik/certs/local-key.pem
  stores:
    default:
      defaultCertificate:
        certFile: /etc/traefik/certs/local-cert.pem
        keyFile: /etc/traefik/certs/local-key.pem
`

func SetupTraefik() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	traefikDir := filepath.Join(homeDir, ".deck-traefik")

	// Create Traefik directory structure
	dirs := []string{
		traefikDir,
		filepath.Join(traefikDir, "certs"),
		filepath.Join(traefikDir, "dynamic"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Generate docker-compose.yml
	composePath := filepath.Join(traefikDir, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte(traefikDockerCompose), 0644); err != nil {
		return fmt.Errorf("failed to create traefik docker-compose.yml: %w", err)
	}

	// Generate dynamic config
	dynamicPath := filepath.Join(traefikDir, "dynamic", "tls.yml")
	if err := os.WriteFile(dynamicPath, []byte(traefikDynamicConfig), 0644); err != nil {
		return fmt.Errorf("failed to create traefik dynamic config: %w", err)
	}

	// Generate SSL certificates
	if err := generateSSLCerts(filepath.Join(traefikDir, "certs")); err != nil {
		return fmt.Errorf("failed to generate SSL certificates: %w", err)
	}

	// Start Traefik
	cmd := exec.Command("docker", "compose", "up", "-d")
	cmd.Dir = traefikDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start Traefik: %w", err)
	}

	return nil
}

func IsTraefikRunning() bool {
	cmd := exec.Command("docker", "ps", "--filter", "name=deck_traefik", "--format", "{{.Names}}")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return len(output) > 0
}

func generateSSLCerts(certsDir string) error {
	certPath := filepath.Join(certsDir, "local-cert.pem")
	keyPath := filepath.Join(certsDir, "local-key.pem")

	// Check if certificates already exist
	if _, err := os.Stat(certPath); err == nil {
		return nil
	}

	// Generate self-signed certificate for *.test domain
	cmd := exec.Command("openssl", "req", "-x509", "-newkey", "rsa:4096",
		"-keyout", keyPath,
		"-out", certPath,
		"-days", "365",
		"-nodes",
		"-subj", "/CN=*.test/O=Deck Local Development",
		"-addext", "subjectAltName=DNS:*.test,DNS:test")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate SSL certificate: %w", err)
	}

	return nil
}

func GetTraefikDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ".deck-traefik"), nil
}
