# Sistema de Versionamento do Deck

Este documento descreve como o Deck gerencia versões do Magento e suas dependências.

## Arquitetura

O sistema utiliza arquivos YAML individuais para cada versão do Magento, permitindo:
- ✅ Fácil manutenção e atualização
- ✅ Adição de novas versões sem recompilação do código
- ✅ Controle granular por versão (incluindo patches, betas e alphas)
- ✅ Fallback inteligente para versões base

## Estrutura de Arquivos

```
internal/magento/
├── versions.go              # Lógica de carregamento e busca
└── versions/                # Diretório de versões
    ├── README.md           # Documentação
    ├── 2.4.7.yaml          # Versão base
    ├── 2.4.7-p1.yaml       # Patch 1
    ├── 2.4.7-p2.yaml       # Patch 2
    ├── 2.4.8.yaml
    ├── 2.4.8-p3.yaml
    ├── 2.4.9-alpha3.yaml   # Alpha
    └── 2.4.9-beta1.yaml    # Beta
```

## Formato do Arquivo YAML

```yaml
version: 2.4.8-p3
php: 8.3
nginx: 1.28
mariadb: 11.4
opensearch: 3
redis: 7.4
rabbitmq: 4.1
```

## Como Funciona

### 1. Carregamento (Embed FS)

Os arquivos YAML são embutidos no binário usando `go:embed`:

```go
//go:embed versions/*.yaml
var versionsFS embed.FS
```

Isso significa:
- Não precisa distribuir arquivos separados
- Binário único e portável
- Carregamento rápido na inicialização

### 2. Cache em Memória

Todas as versões são carregadas na inicialização e armazenadas em um map:

```go
versionCache map[string]*MagentoRequirements
```

### 3. Busca com Fallback

Quando você especifica `magento: 2.4.8-p5`:

1. Busca exata: `2.4.8-p5.yaml`
2. Se não encontrar, busca versão base: `2.4.8.yaml`
3. Se não encontrar, retorna erro com versões disponíveis

### 4. Override de Versões

```yaml
magento: 2.4.8-p3  # Define as versões padrão
php: 8.4           # Mas você pode sobrescrever individualmente
```

## Adicionando Novas Versões

### Opção 1: Script (Recomendado)

```bash
./scripts/add-magento-version.sh 2.4.9
```

### Opção 2: Manual

1. Crie `internal/magento/versions/2.4.9.yaml`
2. Preencha com as versões compatíveis
3. Recompile: `make build && make install`

## Convenções de Nomenclatura

| Tipo | Formato | Exemplo |
|------|---------|---------|
| Release | X.Y.Z | `2.4.8` |
| Patch | X.Y.Z-pN | `2.4.8-p1` |
| Beta | X.Y.Z-betaN | `2.4.9-beta1` |
| Alpha | X.Y.Z-alphaN | `2.4.9-alpha3` |
| RC | X.Y.Z-rcN | `2.4.9-rc1` |

## Funções Disponíveis

### `GetRequirements(version string)`
Retorna os requisitos para uma versão específica.

```go
req, err := magento.GetRequirements("2.4.8-p3")
// req.PHP = "8.3"
// req.Nginx = "1.28"
```

### `GetSupportedVersions()`
Lista todas as versões suportadas.

```go
versions := magento.GetSupportedVersions()
// ["2.4.7", "2.4.7-p1", "2.4.8", ...]
```

### `GetLatestVersion()`
Retorna a versão mais recente disponível.

```go
latest := magento.GetLatestVersion()
// "2.4.9-beta1"
```

### `ListVersionsByMajorMinor()`
Agrupa versões por série.

```go
grouped := magento.ListVersionsByMajorMinor()
// {
//   "2.4": ["2.4.7", "2.4.7-p1", "2.4.8", "2.4.8-p3"],
// }
```

## Fontes de Referência

As versões são baseadas nos requisitos oficiais:

- [Magento System Requirements](https://experienceleague.adobe.com/docs/commerce-operations/installation-guide/system-requirements.html)
- [Magento Release Notes](https://experienceleague.adobe.com/docs/commerce-operations/release/notes/overview.html)

## Roadmap

Melhorias planejadas:

- [ ] Validação automática de compatibilidade
- [ ] Comando `deck versions` para listar versões
- [ ] Comando `deck versions --latest` para mostrar a mais recente
- [ ] Warning se usar versões incompatíveis
- [ ] Download automático de novas versões de um repositório
- [ ] Suporte para Elasticsearch além de OpenSearch
- [ ] Versões LTS marcadas

## Contribuindo

Para adicionar suporte a uma nova versão do Magento:

1. Consulte os [requisitos oficiais](https://experienceleague.adobe.com/docs/commerce-operations/installation-guide/system-requirements.html)
2. Crie o arquivo YAML com as versões compatíveis
3. Teste localmente
4. Envie um Pull Request

## Changelog de Versões

### v1.0.0
- Sistema de versionamento baseado em YAML
- Suporte para Magento 2.4.7+
- Fallback inteligente para versões base
- Script helper para adicionar versões
