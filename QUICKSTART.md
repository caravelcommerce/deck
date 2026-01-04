# Guia Rápido - Deck

## Instalação

### Instalação com um comando (Recomendado)

```bash
curl -fsSL https://raw.githubusercontent.com/caravelcommerce/deck/main/scripts/install.sh | bash
```

Depois recarregue seu shell:

```bash
source ~/.bashrc  # ou ~/.zshrc se você usa zsh
```

## Configuração Básica

1. **Crie um `deck.yaml` na raiz do seu projeto Magento:**

```yaml
name: meu-projeto
magento: 2.4.8-p3
```

2. **Configure o ambiente:**

```bash
deck setup
```

3. **Inicie os containers:**

```bash
deck start
```

4. **Acesse seu site:**

```
https://meu-projeto.test
```

## Comandos Úteis

```bash
# Executar comandos Magento
deck bin/magento setup:upgrade
deck bin/magento cache:flush
deck bin/magento indexer:reindex

# Parar o ambiente
deck stop

# Reiniciar o ambiente
deck stop && deck start
```

## Exemplos de Configuração

### Básico (Recomendado)
```yaml
name: loja
magento: 2.4.8-p3
```

### Com OpenSwoole
```yaml
name: loja
magento: 2.4.8-p3
openswoole: true
```

### Versão Customizada do PHP
```yaml
name: loja
magento: 2.4.8-p3
php: 8.4  # Sobrescreve o padrão (8.3)
```

### Totalmente Manual
```yaml
name: loja
php: 8.3
nginx: 1.28
mariadb: 11.4
opensearch: 3
redis: 7.4
rabbitmq: 4.1
```

## Múltiplos Projetos

Você pode rodar vários projetos ao mesmo tempo:

```bash
# Projeto 1
cd ~/projetos/loja-a
deck start

# Projeto 2
cd ~/projetos/loja-b
deck start
```

Acesse:
- https://loja-a.test
- https://loja-b.test

## Solução Rápida de Problemas

### SSL não funciona
```bash
# macOS
sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain ~/.deck-traefik/certs/local-cert.pem

# Linux
sudo cp ~/.deck-traefik/certs/local-cert.pem /usr/local/share/ca-certificates/deck-local.crt
sudo update-ca-certificates
```

### Container não inicia
```bash
cd .deck
docker compose logs
docker compose down
docker compose up -d --build
```

### Porta ocupada
```bash
# macOS
sudo lsof -i :80
sudo lsof -i :443

# Linux
sudo netstat -tulpn | grep :80
sudo netstat -tulpn | grep :443
```

## Versões Magento Suportadas

| Versão  | Status    |
|---------|-----------|
| 2.4.8   | ✅ Testado |
| 2.4.7   | ✅ Testado |
| 2.4.6   | ✅ Testado |
| 2.4.5   | ✅ Testado |
| 2.4.4   | ✅ Testado |
| 2.4.3   | ✅ Testado |

## Mais Informações

Consulte o [README.md](README.md) para documentação completa.
