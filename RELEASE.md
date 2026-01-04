# Como Criar um Release

Este documento explica como criar uma nova versão do Deck e disponibilizá-la para instalação via script.

## Pré-requisitos

1. Ter permissões de push no repositório
2. Ter o Git configurado corretamente
3. Código testado e funcionando

## Processo de Release

### 1. Certifique-se que tudo está commitado

```bash
git status
git add .
git commit -m "Prepare for release vX.Y.Z"
```

### 2. Crie uma tag de versão

A tag deve seguir o padrão semântico versionamento (semver):
- **v1.0.0** - Release inicial
- **v1.0.1** - Bug fixes
- **v1.1.0** - Novas features (compatível com versão anterior)
- **v2.0.0** - Breaking changes

```bash
# Exemplo para versão 1.0.0
git tag -a v1.0.0 -m "Release v1.0.0 - Initial release"
```

### 3. Faça push da tag

```bash
git push origin v1.0.0
```

### 4. Processo automático

Quando você faz push da tag, o GitHub Actions automaticamente:

1. Faz build dos binários para:
   - Linux AMD64
   - Linux ARM64
   - macOS AMD64 (Intel)
   - macOS ARM64 (Apple Silicon)

2. Gera checksums SHA256 para cada binário

3. Cria um release no GitHub com:
   - Todos os binários anexados
   - Checksums
   - Instruções de instalação
   - Link para o CHANGELOG

4. Publica o release

### 5. Verificar o release

1. Acesse: https://github.com/caravelcommerce/deck/releases
2. Confirme que o release v1.0.0 foi criado
3. Verifique se todos os 4 binários estão anexados
4. Teste a instalação:

```bash
curl -fsSL https://raw.githubusercontent.com/caravelcommerce/deck/main/scripts/install.sh | bash
```

## Estrutura de Versionamento

### Versão Major (X.0.0)
- Mudanças incompatíveis com versões anteriores
- Remoção de features
- Alterações significativas na API/CLI

### Versão Minor (1.X.0)
- Novas features compatíveis com versão anterior
- Adição de comandos ou opções
- Melhorias de funcionalidade

### Versão Patch (1.0.X)
- Correções de bugs
- Pequenas melhorias
- Atualizações de documentação

## Atualizar o CHANGELOG

Antes de criar um release, atualize o [CHANGELOG.md](CHANGELOG.md):

```markdown
## [1.0.0] - 2024-01-04

### Added
- Script de instalação com um comando
- Suporte para detecção automática de versão do Magento
- Criação automática de deck.yaml

### Changed
- Melhorias na documentação

### Fixed
- Correção de bugs na geração de configuração
```

## Testando Localmente

Antes de criar o release oficial, você pode testar o build localmente:

```bash
# Linux AMD64
GOOS=linux GOARCH=amd64 go build -o deck-linux-amd64

# macOS ARM64
GOOS=darwin GOARCH=arm64 go build -o deck-darwin-arm64
```

## Troubleshooting

### Release não foi criado
- Verifique se a tag foi enviada: `git ls-remote --tags origin`
- Verifique os logs do GitHub Actions
- Certifique-se que o workflow está ativado

### Build falhou
- Verifique os logs do GitHub Actions
- Teste o build localmente
- Certifique-se que não há erros de compilação

### Script de instalação não funciona
- Verifique se o release está marcado como "Latest"
- Confirme que os binários estão anexados
- Teste manualmente o download dos binários

## Exemplo Completo

```bash
# 1. Finalize suas mudanças
git add .
git commit -m "Add auto-detection of Magento version"

# 2. Push para o repositório
git push origin main

# 3. Crie a tag
git tag -a v1.0.0 -m "Release v1.0.0 - Initial release with installer"

# 4. Push da tag
git push origin v1.0.0

# 5. Aguarde o GitHub Actions terminar (2-3 minutos)

# 6. Teste a instalação
curl -fsSL https://raw.githubusercontent.com/caravelcommerce/deck/main/scripts/install.sh | bash
deck --version
```

## Hotfix Release

Para correções urgentes:

```bash
# Crie um branch de hotfix
git checkout -b hotfix/v1.0.1

# Faça as correções necessárias
# ... edite arquivos ...

# Commit
git commit -am "Fix critical bug in setup command"

# Merge para main
git checkout main
git merge hotfix/v1.0.1

# Crie a tag
git tag -a v1.0.1 -m "Release v1.0.1 - Hotfix for setup command"

# Push tudo
git push origin main
git push origin v1.0.1

# Delete o branch de hotfix
git branch -d hotfix/v1.0.1
```

## Links Úteis

- [Releases do GitHub](https://github.com/caravelcommerce/deck/releases)
- [GitHub Actions](https://github.com/caravelcommerce/deck/actions)
- [Semantic Versioning](https://semver.org/)
