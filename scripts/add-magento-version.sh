#!/bin/bash

# Script para adicionar uma nova versão do Magento
# Uso: ./add-magento-version.sh 2.4.9-p1

set -e

if [ -z "$1" ]; then
    echo "Uso: $0 <versão>"
    echo "Exemplo: $0 2.4.9-p1"
    exit 1
fi

VERSION=$1
VERSIONS_DIR="$(cd "$(dirname "$0")/../internal/magento/versions" && pwd)"
FILE_PATH="$VERSIONS_DIR/$VERSION.yaml"

if [ -f "$FILE_PATH" ]; then
    echo "❌ Arquivo já existe: $FILE_PATH"
    exit 1
fi

echo "📝 Criando arquivo de versão para Magento $VERSION"
echo ""
echo "Por favor, informe as versões dos serviços:"
echo ""

read -p "PHP version (default: 8.3): " PHP_VERSION
PHP_VERSION=${PHP_VERSION:-8.3}

read -p "Nginx version (default: 1.28): " NGINX_VERSION
NGINX_VERSION=${NGINX_VERSION:-1.28}

read -p "MariaDB version (default: 11.4): " MARIADB_VERSION
MARIADB_VERSION=${MARIADB_VERSION:-11.4}

read -p "OpenSearch version (default: 3): " OPENSEARCH_VERSION
OPENSEARCH_VERSION=${OPENSEARCH_VERSION:-3}

read -p "Redis version (default: 7.4): " REDIS_VERSION
REDIS_VERSION=${REDIS_VERSION:-7.4}

read -p "RabbitMQ version (default: 4.1): " RABBITMQ_VERSION
RABBITMQ_VERSION=${RABBITMQ_VERSION:-4.1}

cat > "$FILE_PATH" << EOF
version: $VERSION
php: $PHP_VERSION
nginx: $NGINX_VERSION
mariadb: $MARIADB_VERSION
opensearch: $OPENSEARCH_VERSION
redis: $REDIS_VERSION
rabbitmq: $RABBITMQ_VERSION
EOF

echo ""
echo "✅ Arquivo criado: $FILE_PATH"
echo ""
echo "Conteúdo:"
cat "$FILE_PATH"
echo ""
echo "Para testar, recompile o projeto:"
echo "  make build"
