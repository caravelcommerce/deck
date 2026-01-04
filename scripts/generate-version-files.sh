#!/bin/bash

# Script para gerar arquivos de versão do Magento com estrutura completa
# Uso: ./generate-version-files.sh

set -e

VERSIONS_DIR="$(cd "$(dirname "$0")/../internal/magento/versions" && pwd)"

# Template base
create_version_file() {
    local version=$1
    local php=$2
    local nginx=$3
    local mariadb=$4
    local opensearch=$5
    local redis=$6
    local rabbitmq=$7

    cat > "$VERSIONS_DIR/$version.yaml" << EOF
version: $version

# PHP Configuration
php:
  version: $php
  extensions:
    - bcmath
    - gd
    - intl
    - mbstring
    - pdo_mysql
    - soap
    - sockets
    - xsl
    - zip
    - opcache

# Nginx Configuration
nginx:
  version: $nginx
  configuration:
    client_max_body_size: 64M
    fastcgi_read_timeout: 600
    fastcgi_connect_timeout: 600

# MariaDB Configuration
mariadb:
  version: $mariadb
  configuration:
    max_connections: 500
    innodb_buffer_pool_size: 1G
    innodb_log_file_size: 256M
    max_allowed_packet: 256M

# OpenSearch Configuration
opensearch:
  version: $opensearch
  configuration:
    cluster.name: magento-cluster
    network.host: 0.0.0.0
    discovery.type: single-node
    OPENSEARCH_JAVA_OPTS: "-Xms512m -Xmx512m"
    DISABLE_SECURITY_PLUGIN: "true"

# Redis Configuration
redis:
  version: $redis
  configuration:
    maxmemory: 256mb
    maxmemory-policy: allkeys-lru
    appendonly: "yes"

# RabbitMQ Configuration
rabbitmq:
  version: $rabbitmq
  configuration:
    RABBITMQ_DEFAULT_USER: guest
    RABBITMQ_DEFAULT_PASS: guest
EOF

    echo "✅ Created $version.yaml"
}

echo "📝 Generating Magento version files..."
echo ""

# Magento 2.4.7 series
create_version_file "2.4.7"    "8.3" "1.28" "11.4" "2.12" "7.4" "3.13"
create_version_file "2.4.7-p1" "8.3" "1.28" "11.4" "2.12" "7.4" "3.13"
create_version_file "2.4.7-p2" "8.3" "1.28" "11.4" "2.12" "7.4" "3.13"

# Magento 2.4.8 series
create_version_file "2.4.8"    "8.3" "1.28" "11.4" "3" "7.4" "4.1"
create_version_file "2.4.8-p1" "8.3" "1.28" "11.4" "3" "7.4" "4.1"
create_version_file "2.4.8-p2" "8.3" "1.28" "11.4" "3" "7.4" "4.1"
create_version_file "2.4.8-p3" "8.3" "1.28" "11.4" "3" "7.4" "4.1"

# Magento 2.4.9 series (development)
create_version_file "2.4.9-alpha3" "8.3" "1.28" "11.4" "3" "7.4" "4.1"
create_version_file "2.4.9-beta1"  "8.3" "1.28" "11.4" "3" "7.4" "4.1"

echo ""
echo "✨ All version files generated successfully!"
echo "📁 Location: $VERSIONS_DIR"
