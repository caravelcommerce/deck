# Magento Version Compatibility Files

Este diretório contém os arquivos YAML com as definições de compatibilidade de software para cada versão do Magento 2.

## Estrutura dos Arquivos

Cada arquivo YAML representa uma versão específica do Magento e define as versões compatíveis dos serviços necessários.

### Formato do Arquivo

```yaml
version: 2.4.8-p3
php: 8.3
nginx: 1.28
mariadb: 11.4
opensearch: 3
redis: 7.4
rabbitmq: 4.1
```

### Convenção de Nomenclatura

- Nome do arquivo: `{versão}.yaml`
- Exemplos:
  - `2.4.8.yaml` - Versão principal
  - `2.4.8-p1.yaml` - Patch release
  - `2.4.9-beta1.yaml` - Beta release
  - `2.4.9-alpha3.yaml` - Alpha release

## Como Adicionar uma Nova Versão

1. Crie um novo arquivo YAML com o nome da versão
2. Preencha as versões compatíveis baseando-se nos [requisitos oficiais do Magento](https://experienceleague.adobe.com/docs/commerce-operations/installation-guide/system-requirements.html)
3. O arquivo será automaticamente carregado na próxima compilação

### Exemplo: Adicionando Magento 2.4.9

```bash
# Crie o arquivo
cat > 2.4.9.yaml << EOF
version: 2.4.9
php: 8.4
nginx: 1.28
mariadb: 11.4
opensearch: 3
redis: 7.4
rabbitmq: 4.1
EOF
```

## Fallback de Versões

O sistema possui fallback inteligente:

- Se `2.4.8-p5` não existir, ele usa `2.4.8`
- Se uma versão específica não for encontrada, retorna erro com lista de versões disponíveis

## Versões Atualmente Suportadas

### Magento 2.4.7
- 2.4.7
- 2.4.7-p1
- 2.4.7-p2
- 2.4.7-p3

### Magento 2.4.8
- 2.4.8
- 2.4.8-p1
- 2.4.8-p2
- 2.4.8-p3

### Magento 2.4.9 (Desenvolvimento)
- 2.4.9-alpha3
- 2.4.9-beta1

## Fontes de Referência

- [Magento System Requirements](https://experienceleague.adobe.com/docs/commerce-operations/installation-guide/system-requirements.html)
- [Magento Release Notes](https://experienceleague.adobe.com/docs/commerce-operations/release/notes/overview.html)
- [Magento DevDocs](https://developer.adobe.com/commerce/php/development/)

## Atualizações

Para manter as versões atualizadas, consulte regularmente:
- Adobe Commerce Release Schedule
- Magento Open Source GitHub Releases
- Magento DevBlog
