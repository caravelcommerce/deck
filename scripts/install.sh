#!/bin/bash
set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configurações
REPO="caravelcommerce/deck"
INSTALL_DIR="$HOME/.deck"
BIN_DIR="$HOME/.local/bin"

# Funções auxiliares
print_info() {
    echo -e "${BLUE}ℹ${NC} $1"
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

# Detecta o sistema operacional e arquitetura
detect_platform() {
    OS="$(uname -s)"
    ARCH="$(uname -m)"

    case "$OS" in
        Linux*)
            PLATFORM="linux"
            ;;
        Darwin*)
            PLATFORM="darwin"
            ;;
        *)
            print_error "Sistema operacional não suportado: $OS"
            exit 1
            ;;
    esac

    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            print_error "Arquitetura não suportada: $ARCH"
            exit 1
            ;;
    esac

    print_info "Plataforma detectada: $PLATFORM/$ARCH"
}

# Obtém a última versão do GitHub
get_latest_version() {
    print_info "Buscando última versão..."

    VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        print_error "Não foi possível obter a última versão"
        exit 1
    fi

    print_success "Última versão: $VERSION"
}

# Download do binário
download_deck() {
    BINARY_NAME="deck-${PLATFORM}-${ARCH}"
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY_NAME"

    print_info "Baixando Deck $VERSION..."

    # Cria diretório temporário
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"

    # Faz o download
    if ! curl -fsSL "$DOWNLOAD_URL" -o deck; then
        print_error "Falha ao baixar $DOWNLOAD_URL"
        rm -rf "$TMP_DIR"
        exit 1
    fi

    # Torna o binário executável
    chmod +x deck

    print_success "Download concluído"
}

# Instala o binário
install_deck() {
    print_info "Instalando Deck..."

    # Cria diretório de instalação se não existir
    mkdir -p "$BIN_DIR"

    # Move o binário
    mv deck "$BIN_DIR/deck"

    # Remove diretório temporário
    cd ~
    rm -rf "$TMP_DIR"

    print_success "Deck instalado em $BIN_DIR/deck"
}

# Configura o PATH
setup_path() {
    # Detecta o shell
    SHELL_NAME=$(basename "$SHELL")

    case "$SHELL_NAME" in
        bash)
            PROFILE="$HOME/.bashrc"
            if [ "$(uname)" = "Darwin" ]; then
                PROFILE="$HOME/.bash_profile"
            fi
            ;;
        zsh)
            PROFILE="$HOME/.zshrc"
            ;;
        fish)
            PROFILE="$HOME/.config/fish/config.fish"
            ;;
        *)
            PROFILE="$HOME/.profile"
            ;;
    esac

    # Verifica se o BIN_DIR já está no PATH
    if [[ ":$PATH:" == *":$BIN_DIR:"* ]]; then
        print_success "PATH já configurado"
        return
    fi

    # Adiciona ao PATH no arquivo de perfil
    if [ "$SHELL_NAME" = "fish" ]; then
        echo "" >> "$PROFILE"
        echo "# Deck CLI" >> "$PROFILE"
        echo "set -gx PATH \$PATH $BIN_DIR" >> "$PROFILE"
    else
        echo "" >> "$PROFILE"
        echo "# Deck CLI" >> "$PROFILE"
        echo "export PATH=\"\$PATH:$BIN_DIR\"" >> "$PROFILE"
    fi

    print_success "PATH configurado em $PROFILE"
    print_warning "Execute 'source $PROFILE' ou abra um novo terminal para aplicar as mudanças"
}

# Verifica a instalação
verify_installation() {
    print_info "Verificando instalação..."

    if [ -x "$BIN_DIR/deck" ]; then
        VERSION_OUTPUT=$("$BIN_DIR/deck" --version 2>&1 || echo "")
        print_success "Deck instalado com sucesso!"
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo -e "${GREEN}  Deck - Magento 2 Development Environment${NC}"
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
        echo -e "  Versão: ${BLUE}$VERSION${NC}"
        echo -e "  Local: ${BLUE}$BIN_DIR/deck${NC}"
        echo ""
        echo -e "${YELLOW}Próximos passos:${NC}"
        echo ""
        echo "  1. Recarregue seu shell:"
        echo -e "     ${BLUE}source $PROFILE${NC}"
        echo ""
        echo "  2. Navegue até seu projeto Magento:"
        echo -e "     ${BLUE}cd /path/to/magento${NC}"
        echo ""
        echo "  3. Execute o setup:"
        echo -e "     ${BLUE}deck setup${NC}"
        echo ""
        echo "  4. Inicie o ambiente:"
        echo -e "     ${BLUE}deck start${NC}"
        echo ""
        echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
        echo ""
    else
        print_error "Falha na verificação da instalação"
        exit 1
    fi
}

# Função principal
main() {
    echo ""
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${BLUE}  Instalador do Deck${NC}"
    echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""

    detect_platform
    get_latest_version
    download_deck
    install_deck
    setup_path
    verify_installation
}

# Executa
main
