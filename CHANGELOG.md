# Changelog

## [1.0.0] - 2026-01-04

### Adicionado
- Sistema completo de desenvolvimento Docker para Magento 2
- CLI com comandos: setup, start, stop, bin/magento
- Suporte para múltiplos projetos simultâneos
- Traefik como reverse proxy com SSL automático
- Domínios `.test` automáticos com certificados SSL
- Configuração personalizável via `deck.yaml`
- **Auto-detecção de versões baseada na versão do Magento**
  - Especifique apenas `magento: 2.4.8-p3` e o Deck configura automaticamente
  - Matriz de compatibilidade para Magento 2.4.3 até 2.4.8
  - Suporte para sobrescrever versões individuais
  - Detecção inteligente de versões mais próximas
- Suporte para OpenSwoole (opcional)
- Templates Docker otimizados para:
  - Nginx (1.22 - 1.28)
  - PHP (8.1 - 8.4) com extensões Magento
  - MariaDB (10.4 - 11.4)
  - OpenSearch (1.2 - 3)
  - Redis (6.2 - 7.4)
  - RabbitMQ (3.9 - 4.1)

### Funcionalidades
- `deck setup` - Configura ambiente na pasta `.deck`
  - Exibe versões detectadas/configuradas
  - Mostra se está usando auto-detecção ou configuração manual
- `deck start` - Inicia containers Docker
- `deck stop` - Para containers Docker
- `deck bin/magento` - Executa comandos Magento CLI

### Extensões PHP Opcionais
- **OpenSwoole** - Programação assíncrona e corrotinas
  - Configurável via `openswoole: true` no deck.yaml
  - Instalação automática durante o build
  - Configuração otimizada com `use_shortname = 'Off'`

### Matriz de Compatibilidade
| Magento | PHP | Nginx | MariaDB | OpenSearch | Redis | RabbitMQ |
|---------|-----|-------|---------|------------|-------|----------|
| 2.4.8   | 8.3 | 1.28  | 11.4    | 3          | 7.4   | 4.1      |
| 2.4.7   | 8.3 | 1.28  | 11.4    | 2.12       | 7.4   | 3.13     |
| 2.4.6   | 8.2 | 1.24  | 10.6    | 2.12       | 7.2   | 3.13     |
| 2.4.5   | 8.1 | 1.24  | 10.6    | 2.5        | 7.0   | 3.11     |
| 2.4.4   | 8.1 | 1.22  | 10.6    | 1.2        | 7.0   | 3.9      |
| 2.4.3   | 8.1 | 1.22  | 10.4    | 1.2        | 6.2   | 3.9      |
