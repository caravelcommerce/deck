# Deck - Magento 2 Docker Development Environment

Um sistema completo de desenvolvimento para Magento 2 baseado em Docker com suporte para múltiplos projetos simultâneos.

## Características

- ✨ Suporte para múltiplos projetos Magento rodando simultaneamente
- 🔧 Configuração personalizável por projeto via `deck.yaml`
- 🌐 Domínios automáticos `.test` com SSL
- 🚀 CLI simples e intuitivo
- 🔒 Traefik como reverse proxy com SSL automático
- 📦 Containers otimizados para Magento 2

## Pré-requisitos

- Docker (>= 24.0)
- Docker Compose (>= 2.20)

### macOS
```bash
brew install docker docker-compose
```

### Linux (Ubuntu/Debian)
```bash
sudo apt update
sudo apt install docker.io docker-compose
```

## Instalação

### Instalação para Repositório Privado

Como este é um repositório privado, você precisa de um GitHub Personal Access Token para instalar.

#### 1. Crie um Personal Access Token

1. Acesse: https://github.com/settings/tokens/new
2. Dê um nome (ex: "Deck CLI")
3. Selecione a permissão: **repo** (Full control of private repositories)
4. Clique em "Generate token"
5. Copie o token gerado

#### 2. Instale o Deck

```bash
# Defina o token (temporário - apenas para esta sessão)
export GITHUB_TOKEN=seu_token_aqui

# Execute o instalador
curl -fsSL https://raw.githubusercontent.com/caravelcommerce/deck/main/scripts/install-private.sh | bash
```

**OU** defina o token permanentemente:

```bash
# Adicione ao seu .bashrc ou .zshrc
echo 'export GITHUB_TOKEN=seu_token_aqui' >> ~/.bashrc
source ~/.bashrc

# Execute o instalador
curl -fsSL https://raw.githubusercontent.com/caravelcommerce/deck/main/scripts/install-private.sh | bash
```

Após a instalação, recarregue seu shell:

```bash
source ~/.bashrc  # ou ~/.zshrc se você usa zsh
```

### Instalação Manual

Se preferir instalar manualmente:

1. Baixe o binário apropriado para sua plataforma na [página de releases](https://github.com/caravelcommerce/deck/releases/latest):
   - **Linux AMD64**: `deck-linux-amd64`
   - **Linux ARM64**: `deck-linux-arm64`
   - **macOS Intel**: `deck-darwin-amd64`
   - **macOS Apple Silicon**: `deck-darwin-arm64`

2. Torne o binário executável e mova para seu PATH:

```bash
chmod +x deck-*
mkdir -p ~/.local/bin
mv deck-* ~/.local/bin/deck
```

3. Adicione ao PATH (se necessário):

```bash
# Para bash
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.bashrc

# Para zsh
echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
```

### Verificar Instalação

```bash
deck --version
```

### Atualização

Para atualizar para a versão mais recente:

```bash
curl -fsSL https://raw.githubusercontent.com/caravelcommerce/deck/main/scripts/install.sh | bash
```

O script detectará a instalação existente e substituirá pelo novo binário.

## Configuração do Projeto

### 1. Crie um arquivo `deck.yaml` na raiz do seu projeto Magento

#### Opção 1: Auto-detecção baseada na versão do Magento (Recomendado)

```yaml
name: demo
magento: 2.4.8-p3  # O Deck detecta automaticamente as versões compatíveis

# Extensões PHP opcionais
openswoole: false
```

**Versões Magento suportadas:**
- Série 2.4.7: `2.4.7`, `2.4.7-p1`, `2.4.7-p2`, `2.4.7-p3`
- Série 2.4.8: `2.4.8`, `2.4.8-p1`, `2.4.8-p2`, `2.4.8-p3`
- Série 2.4.9: `2.4.9-alpha3`, `2.4.9-beta1`

> **Nota:** Novas versões são adicionadas regularmente. Consulte [internal/magento/versions/](internal/magento/versions/) para a lista completa.

#### Opção 2: Especificar versões manualmente

```yaml
name: demo
nginx: 1.28
php: 8.3
mariadb: 11.4
opensearch: 3
redis: 7.4
rabbitmq: 4.1
openswoole: false
```

#### Opção 3: Híbrido (Auto-detecção + Override)

```yaml
name: demo
magento: 2.4.8-p3  # Usa versões compatíveis com Magento 2.4.8
php: 8.4           # Mas sobrescreve o PHP para 8.4
openswoole: true   # E habilita OpenSwoole
```

### 2. Execute o setup
```bash
deck setup
```

Este comando irá:
- Criar a pasta `.deck` com todas as configurações Docker
- Configurar o Traefik reverse proxy (se ainda não estiver rodando)
- Gerar certificados SSL para `*.test`
- Adicionar `.deck/` ao `.gitignore`

### 3. Inicie o ambiente
```bash
deck start
```

Seu projeto estará disponível em `https://{name}.test` (exemplo: `https://demo.test`)

## Comandos Disponíveis

### `deck setup`
Configura o ambiente Docker baseado no arquivo `deck.yaml`. Cria a pasta `.deck` com todos os arquivos necessários.

```bash
deck setup
```

### `deck start`
Inicia todos os containers Docker do projeto.

```bash
deck start
```

### `deck stop`
Para todos os containers Docker do projeto.

```bash
deck stop
```

### `deck bin/magento`
Executa comandos do Magento CLI dentro do container PHP.

```bash
# Exemplos:
deck bin/magento setup:upgrade
deck bin/magento cache:flush
deck bin/magento indexer:reindex
deck bin/magento deploy:mode:set developer
```

## Matriz de Compatibilidade Magento

O Deck inclui uma matriz de compatibilidade baseada nos [requisitos oficiais do Magento](https://experienceleague.adobe.com/docs/commerce-operations/installation-guide/system-requirements.html):

| Magento  | PHP  | Nginx | MariaDB | OpenSearch | Redis | RabbitMQ |
|----------|------|-------|---------|------------|-------|----------|
| 2.4.8    | 8.3  | 1.28  | 11.4    | 3          | 7.4   | 4.1      |
| 2.4.7    | 8.3  | 1.28  | 11.4    | 2.12       | 7.4   | 3.13     |
| 2.4.6    | 8.2  | 1.24  | 10.6    | 2.12       | 7.2   | 3.13     |
| 2.4.5    | 8.1  | 1.24  | 10.6    | 2.5        | 7.0   | 3.11     |
| 2.4.4    | 8.1  | 1.22  | 10.6    | 1.2        | 7.0   | 3.9      |
| 2.4.3    | 8.1  | 1.22  | 10.4    | 1.2        | 6.2   | 3.9      |

Quando você especifica `magento: 2.4.8-p3`, o Deck automaticamente usa as versões acima. Você pode sobrescrever qualquer versão individualmente.

## Estrutura de Serviços

Cada projeto Magento terá os seguintes serviços:

- **Nginx** - Servidor web
- **PHP-FPM** - PHP com extensões Magento
- **MariaDB** - Banco de dados
- **OpenSearch** - Motor de busca
- **Redis** - Cache e sessões
- **RabbitMQ** - Fila de mensagens

## Acessando os Serviços

### Web
- URL: `https://{name}.test`

### Banco de Dados
- Host: `{name}_mariadb`
- Port: `3306`
- Database: `magento`
- User: `magento`
- Password: `magento`
- Root Password: `root`

### Redis
- Host: `{name}_redis`
- Port: `6379`

### OpenSearch
- Host: `{name}_opensearch`
- Port: `9200`

### RabbitMQ
- Host: `{name}_rabbitmq`
- Port: `5672`
- Management UI: `http://localhost:15672`
- User: `guest`
- Password: `guest`

### Traefik Dashboard
- URL: `http://localhost:8080`

## Extensões PHP Opcionais

### OpenSwoole

OpenSwoole é uma extensão PHP de alto desempenho que habilita programação assíncrona, corrotinas e suporte nativo para HTTP/WebSocket.

#### Instalação Básica

Para habilitar OpenSwoole no seu projeto:

1. Edite o arquivo `deck.yaml` e defina `openswoole: true`:

```yaml
name: demo
magento: 2.4.8-p3
openswoole: true
```

2. Execute `deck setup` para regenerar as configurações Docker
3. Execute `deck start` para iniciar o ambiente

**Nota:** Quando OpenSwoole está habilitado, a configuração `openswoole.use_shortname = 'Off'` é automaticamente aplicada para evitar conflitos com funções nativas do PHP.

Para verificar se OpenSwoole foi instalado corretamente:

```bash
deck bin/magento exec php -m | grep openswoole
```

#### Swoole HTTP Server em Subdomínio (API Assíncrona)

Você pode expor um servidor Swoole HTTP em um subdomínio separado, perfeito para APIs assíncronas rodando em paralelo com o Magento:

```yaml
name: demo
magento: 2.4.8-p3
openswoole: true
swoole_port: 9501  # Porta onde o Swoole HTTP Server irá rodar
```

Com essa configuração:
- **Magento tradicional**: `https://demo.test` (Nginx + PHP-FPM)
- **Swoole API**: `https://api.demo.test` (Swoole HTTP Server)

**Como usar:**

1. Configure seu módulo Magento para iniciar o Swoole HTTP Server:

```bash
# Inicie o servidor Swoole via Magento CLI
deck bin/magento swoole:server:start
```

2. O servidor Swoole rodando na porta `9501` será automaticamente exposto em:
   - `https://api.demo.test` (via Traefik)
   - Certificado SSL automático

3. Faça requisições para sua API:

```bash
curl https://api.demo.test/api/v1/products
```

**Exemplo de Módulo Magento com Swoole:**

```php
<?php
// app/code/Vendor/Swoole/Console/Command/ServerStart.php
namespace Vendor\Swoole\Console\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;

class ServerStart extends Command
{
    protected function configure()
    {
        $this->setName('swoole:server:start')
             ->setDescription('Start Swoole HTTP Server');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        $server = new \Swoole\HTTP\Server("0.0.0.0", 9501);

        $server->on("request", function ($request, $response) {
            $response->header("Content-Type", "application/json");
            $response->end(json_encode([
                'status' => 'success',
                'message' => 'Swoole API is running',
                'path' => $request->server['request_uri']
            ]));
        });

        $output->writeln("Swoole HTTP Server started on port 9501");
        $output->writeln("Access at: https://api.demo.test");

        $server->start();

        return Command::SUCCESS;
    }
}
```

**Casos de Uso:**
- APIs REST assíncronas de alta performance
- WebSocket servers para comunicação em tempo real
- Filas de processamento assíncronas
- Microserviços isolados do Magento tradicional

## Configuração SSL

Os certificados SSL são gerados automaticamente durante o `deck setup`. Para confiar no certificado:

### macOS
```bash
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.deck-traefik/certs/local-cert.pem
```

### Linux
```bash
sudo cp ~/.deck-traefik/certs/local-cert.pem /usr/local/share/ca-certificates/deck-local.crt
sudo update-ca-certificates
```

## Múltiplos Projetos

Para rodar múltiplos projetos Magento simultaneamente:

1. Cada projeto deve ter seu próprio `deck.yaml` com um `name` único
2. Execute `deck setup` e `deck start` em cada projeto
3. Todos os projetos compartilham o mesmo Traefik reverse proxy

Exemplo:
```
~/projetos/loja1/deck.yaml  → name: loja1 → https://loja1.test
~/projetos/loja2/deck.yaml  → name: loja2 → https://loja2.test
~/projetos/demo/deck.yaml   → name: demo  → https://demo.test
```

## Adicionando Novas Versões do Magento

O Deck usa arquivos YAML para definir as compatibilidades de cada versão do Magento. Para adicionar uma nova versão:

### Método 1: Script Automatizado

```bash
./scripts/add-magento-version.sh 2.4.9
```

O script irá solicitar as versões de cada serviço e criar o arquivo automaticamente.

### Método 2: Manual

1. Crie um arquivo em `internal/magento/versions/{versão}.yaml`:

```yaml
version: 2.4.9
php: 8.4
nginx: 1.28
mariadb: 11.4
opensearch: 3
redis: 7.4
rabbitmq: 4.1
```

2. Recompile o projeto:

```bash
make build
make install
```

### Estrutura dos Arquivos de Versão

Todos os arquivos de versão ficam em `internal/magento/versions/`:

```
internal/magento/versions/
├── README.md           # Documentação
├── 2.4.7.yaml
├── 2.4.7-p1.yaml
├── 2.4.8.yaml
├── 2.4.8-p3.yaml
└── 2.4.9-beta1.yaml
```

Consulte o [README de versões](internal/magento/versions/README.md) para mais detalhes.

## Solução de Problemas

### Container não inicia
```bash
# Verifique os logs
cd .deck
docker compose logs

# Recrie os containers
docker compose down
docker compose up -d --build
```

### Porta já em uso
Se a porta 80 ou 443 já estiver em uso, pare o serviço conflitante:

```bash
# macOS - Apache
sudo apachectl stop

# Linux - Apache/Nginx
sudo systemctl stop apache2
sudo systemctl stop nginx
```

### SSL não funciona
Certifique-se de que adicionou o certificado aos certificados confiáveis do sistema (veja seção "Configuração SSL").

### Traefik não responde
```bash
# Reinicie o Traefik
cd ~/.deck-traefik
docker compose restart
```

## Desinstalação

Para remover o Deck:

```bash
rm ~/.local/bin/deck
```

Para remover completamente todos os dados:

```bash
# Remove o CLI
rm ~/.local/bin/deck

# Remove a configuração do PATH (edite manualmente ~/.bashrc ou ~/.zshrc)

# Remove o Traefik
cd ~/.deck-traefik
docker compose down -v
rm -rf ~/.deck-traefik

# Em cada projeto, remova os dados Docker
cd seu-projeto
cd .deck
docker compose down -v
rm -rf .deck
```

## Estrutura de Diretórios

```
seu-projeto/
├── deck.yaml           # Configuração do projeto
├── .deck/              # Gerado pelo 'deck setup'
│   ├── docker-compose.yml
│   ├── nginx/
│   │   ├── nginx.conf
│   │   └── default.conf
│   ├── php/
│   │   ├── Dockerfile
│   │   ├── php.ini
│   │   └── php-fpm.conf
│   └── mariadb/
│       └── my.cnf
└── (seus arquivos Magento)
```

## Contribuindo

Contribuições são bem-vindas! Por favor, abra uma issue ou pull request.

## Licença

MIT License

## Suporte

Para reportar problemas ou sugerir melhorias, abra uma issue no GitHub.
