# Guia Swoole - Deck

Este guia mostra como usar OpenSwoole com Deck para criar APIs assíncronas de alta performance rodando em paralelo com o Magento.

## Arquitetura

```
┌─────────────────────────────────────────────────┐
│              Traefik Reverse Proxy              │
│                 (SSL Automático)                │
└─────────────┬──────────────────┬────────────────┘
              │                  │
              ▼                  ▼
    ┌──────────────────┐  ┌──────────────────┐
    │   Nginx + PHP    │  │  Swoole Server   │
    │  (Magento web)   │  │  (API Async)     │
    │ https://demo.test│  │https://api.demo. │
    │                  │  │      test        │
    └──────────────────┘  └──────────────────┘
```

## Configuração Básica

### 1. Configure o deck.yaml

```yaml
name: demo
magento: 2.4.8-p3
openswoole: true
swoole_port: 9501
```

### 2. Execute o Setup

```bash
deck setup
deck start
```

### 3. Resultado

- **Magento Web**: `https://demo.test` (Nginx + PHP-FPM)
- **Swoole API**: `https://api.demo.test` (Swoole HTTP Server na porta 9501)

## Criando um Módulo Magento com Swoole

### Estrutura do Módulo

```
app/code/Vendor/SwooleApi/
├── registration.php
├── etc/
│   ├── module.xml
│   └── di.xml
└── Console/
    └── Command/
        └── ServerStart.php
```

### registration.php

```php
<?php
use Magento\Framework\Component\ComponentRegistrar;

ComponentRegistrar::register(
    ComponentRegistrar::MODULE,
    'Vendor_SwooleApi',
    __DIR__
);
```

### etc/module.xml

```xml
<?xml version="1.0"?>
<config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:noNamespaceSchemaLocation="urn:magento:framework:Module/etc/module.xsd">
    <module name="Vendor_SwooleApi" setup_version="1.0.0"/>
</config>
```

### etc/di.xml

```xml
<?xml version="1.0"?>
<config xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:noNamespaceSchemaLocation="urn:magento:framework:ObjectManager/etc/config.xsd">
    <type name="Magento\Framework\Console\CommandList">
        <arguments>
            <argument name="commands" xsi:type="array">
                <item name="swoole_server_start" xsi:type="object">
                    Vendor\SwooleApi\Console\Command\ServerStart
                </item>
            </argument>
        </arguments>
    </type>
</config>
```

### Console/Command/ServerStart.php

```php
<?php
namespace Vendor\SwooleApi\Console\Command;

use Symfony\Component\Console\Command\Command;
use Symfony\Component\Console\Input\InputInterface;
use Symfony\Component\Console\Output\OutputInterface;
use Magento\Framework\App\ResourceConnection;
use Magento\Framework\App\State;

class ServerStart extends Command
{
    private $resourceConnection;
    private $appState;

    public function __construct(
        ResourceConnection $resourceConnection,
        State $appState,
        string $name = null
    ) {
        $this->resourceConnection = $resourceConnection;
        $this->appState = $appState;
        parent::__construct($name);
    }

    protected function configure()
    {
        $this->setName('swoole:server:start')
             ->setDescription('Start Swoole HTTP Server for API');
    }

    protected function execute(InputInterface $input, OutputInterface $output)
    {
        try {
            // Set area code
            $this->appState->setAreaCode(\Magento\Framework\App\Area::AREA_GLOBAL);
        } catch (\Exception $e) {
            // Already set
        }

        $server = new \Swoole\HTTP\Server("0.0.0.0", 9501);

        // Configurações do servidor
        $server->set([
            'worker_num' => 4,
            'max_request' => 10000,
            'daemonize' => false,
        ]);

        // Handler de requisições
        $server->on("request", function ($request, $response) {
            $this->handleRequest($request, $response);
        });

        $output->writeln("🚀 Swoole HTTP Server started on port 9501");
        $output->writeln("📡 API endpoint: https://api.demo.test");
        $output->writeln("Press Ctrl+C to stop");

        $server->start();

        return Command::SUCCESS;
    }

    private function handleRequest($request, $response)
    {
        // CORS headers
        $response->header("Access-Control-Allow-Origin", "*");
        $response->header("Content-Type", "application/json");

        $path = $request->server['request_uri'];
        $method = $request->server['request_method'];

        // Roteamento básico
        switch ($path) {
            case '/':
                $this->handleRoot($response);
                break;

            case '/api/products':
                $this->handleProducts($response);
                break;

            case '/api/health':
                $this->handleHealth($response);
                break;

            default:
                $this->handleNotFound($response);
                break;
        }
    }

    private function handleRoot($response)
    {
        $response->end(json_encode([
            'status' => 'success',
            'message' => 'Swoole API Server',
            'version' => '1.0.0',
            'endpoints' => [
                '/api/products',
                '/api/health'
            ]
        ], JSON_PRETTY_PRINT));
    }

    private function handleProducts($response)
    {
        // Exemplo: buscar produtos do banco usando Magento
        $connection = $this->resourceConnection->getConnection();
        $table = $this->resourceConnection->getTableName('catalog_product_entity');

        $query = $connection->select()
            ->from($table, ['entity_id', 'sku', 'created_at'])
            ->limit(10);

        $products = $connection->fetchAll($query);

        $response->end(json_encode([
            'status' => 'success',
            'data' => $products,
            'count' => count($products)
        ], JSON_PRETTY_PRINT));
    }

    private function handleHealth($response)
    {
        $response->end(json_encode([
            'status' => 'healthy',
            'timestamp' => time(),
            'server' => 'swoole'
        ]));
    }

    private function handleNotFound($response)
    {
        $response->status(404);
        $response->end(json_encode([
            'status' => 'error',
            'message' => 'Endpoint not found'
        ]));
    }
}
```

## Instalação e Uso

### 1. Instale o módulo

```bash
deck bin/magento module:enable Vendor_SwooleApi
deck bin/magento setup:upgrade
```

### 2. Inicie o servidor Swoole

```bash
deck bin/magento swoole:server:start
```

### 3. Teste a API

```bash
# Health check
curl https://api.demo.test/api/health

# Lista de produtos
curl https://api.demo.test/api/products

# Info do servidor
curl https://api.demo.test/
```

## Casos de Uso

### 1. API REST de Alta Performance

```yaml
# deck.yaml
name: ecommerce-api
magento: 2.4.8-p3
openswoole: true
swoole_port: 9501
```

Endpoints:
- `https://ecommerce-api.test` - Loja Magento tradicional
- `https://api.ecommerce-api.test` - API REST assíncrona

### 2. WebSocket Server

```php
// Modificar ServerStart.php para usar WebSocket
$server = new \Swoole\WebSocket\Server("0.0.0.0", 9501);

$server->on("open", function($server, $request) {
    echo "Connection opened: {$request->fd}\n";
});

$server->on("message", function($server, $frame) {
    $server->push($frame->fd, "Received: {$frame->data}");
});
```

### 3. Microserviço de Processamento

```php
// Processar filas de forma assíncrona
private function handleQueue($response)
{
    go(function() {
        // Corrotina para processar em background
        $this->processQueueItems();
    });

    $response->end(json_encode(['status' => 'processing']));
}
```

## Monitoramento

### Logs do Swoole

```bash
# Ver logs do container PHP (onde Swoole roda)
cd .deck
docker compose logs -f php
```

### Métricas

Adicione endpoints de métricas:

```php
private function handleMetrics($response)
{
    $stats = $server->stats();
    $response->end(json_encode($stats, JSON_PRETTY_PRINT));
}
```

## Performance

### Benchmarks Típicos

- **Nginx + PHP-FPM**: ~1.000 req/s
- **Swoole**: ~10.000+ req/s (10x mais rápido)

### Otimizações

```php
$server->set([
    'worker_num' => 8,              // Número de workers
    'max_request' => 10000,         // Requisições antes de restart
    'task_worker_num' => 4,         // Workers para tarefas assíncronas
    'task_enable_coroutine' => true // Habilita corrotinas
]);
```

## Troubleshooting

### Porta já em uso

```bash
# Verifique se a porta está em uso
lsof -i :9501

# Mude a porta no deck.yaml
swoole_port: 9502
```

### Servidor não responde

```bash
# Verifique se o container está rodando
docker ps | grep php

# Verifique os logs
cd .deck && docker compose logs php
```

### SSL não funciona

Certifique-se de que adicionou o certificado confiável (veja README.md - Configuração SSL).

## Referências

- [OpenSwoole Docs](https://openswoole.com/docs)
- [Swoole HTTP Server](https://openswoole.com/docs/modules/swoole-http-server)
- [Corrotinas](https://openswoole.com/docs/modules/swoole-coroutine)
