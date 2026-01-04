package docker

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"github.com/caravelcommerce/deck/internal/config"
)

const dockerComposeTemplate = `version: '3.8'

services:
  nginx:
    image: nginx:{{.Nginx}}-alpine
    container_name: {{.Name}}_nginx
    volumes:
      - ../:/var/www/html:cached
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/default.conf:/etc/nginx/conf.d/default.conf:ro
    networks:
      - {{.Name}}_network
      - traefik_network
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.{{.Name}}.rule=Host(` + "`{{.Name}}.test`" + `)"
      - "traefik.http.routers.{{.Name}}.entrypoints=websecure"
      - "traefik.http.routers.{{.Name}}.tls=true"
      - "traefik.http.services.{{.Name}}.loadbalancer.server.port=80"
    depends_on:
      - php

  php:
    build:
      context: ./php
      args:
        PHP_VERSION: {{.PHP}}
        INSTALL_OPENSWOOLE: {{.OpenSwoole}}
    container_name: {{.Name}}_php
    volumes:
      - ../:/var/www/html:cached
      - ./php/php.ini:/usr/local/etc/php/php.ini:ro
      - ./php/php-fpm.conf:/usr/local/etc/php-fpm.d/www.conf:ro
    networks:
      - {{.Name}}_network{{if gt .SwoolePort 0}}
      - traefik_network{{end}}
    environment:
      - PHP_IDE_CONFIG=serverName={{.Name}}{{if gt .SwoolePort 0}}
    ports:
      - "{{.SwoolePort}}:{{.SwoolePort}}"
    labels:
      - "traefik.enable=true"
      # Swoole HTTP Server on api subdomain
      - "traefik.http.routers.{{.Name}}-swoole.rule=Host(` + "`api.{{.Name}}.test`" + `)"
      - "traefik.http.routers.{{.Name}}-swoole.entrypoints=websecure"
      - "traefik.http.routers.{{.Name}}-swoole.tls=true"
      - "traefik.http.routers.{{.Name}}-swoole.service={{.Name}}-swoole"
      - "traefik.http.services.{{.Name}}-swoole.loadbalancer.server.port={{.SwoolePort}}"{{end}}
    depends_on:
      - mariadb
      - redis
      - opensearch
      - rabbitmq

  mariadb:
    image: mariadb:{{.MariaDB}}
    container_name: {{.Name}}_mariadb
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: magento
      MYSQL_USER: magento
      MYSQL_PASSWORD: magento
    volumes:
      - mariadb_data:/var/lib/mysql
      - ./mariadb/my.cnf:/etc/mysql/conf.d/custom.cnf:ro
    networks:
      - {{.Name}}_network
    command: --max_allowed_packet=256M

  opensearch:
    image: opensearchproject/opensearch:{{.OpenSearch}}
    container_name: {{.Name}}_opensearch
    environment:
      - discovery.type=single-node
      - "OPENSEARCH_JAVA_OPTS=-Xms512m -Xmx512m"
      - "DISABLE_SECURITY_PLUGIN=true"
    volumes:
      - opensearch_data:/usr/share/opensearch/data
    networks:
      - {{.Name}}_network

  redis:
    image: redis:{{.Redis}}-alpine
    container_name: {{.Name}}_redis
    volumes:
      - redis_data:/data
    networks:
      - {{.Name}}_network

  rabbitmq:
    image: rabbitmq:{{.RabbitMQ}}-management-alpine
    container_name: {{.Name}}_rabbitmq
    environment:
      RABBITMQ_DEFAULT_USER: guest
      RABBITMQ_DEFAULT_PASS: guest
    volumes:
      - rabbitmq_data:/var/lib/rabbitmq
    networks:
      - {{.Name}}_network

networks:
  {{.Name}}_network:
    driver: bridge
  traefik_network:
    external: true

volumes:
  mariadb_data:
  opensearch_data:
  redis_data:
  rabbitmq_data:
`

const nginxConfTemplate = `user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    log_format main '$remote_addr - $remote_user [$time_local] "$request" '
                    '$status $body_bytes_sent "$http_referer" '
                    '"$http_user_agent" "$http_x_forwarded_for"';

    access_log /var/log/nginx/access.log main;

    sendfile on;
    tcp_nopush on;
    tcp_nodelay on;
    keepalive_timeout 65;
    types_hash_max_size 2048;
    client_max_body_size 256M;

    gzip on;
    gzip_disable "msie6";
    gzip_vary on;
    gzip_proxied any;
    gzip_comp_level 6;
    gzip_types text/plain text/css text/xml text/javascript application/json application/javascript application/xml+rss;

    include /etc/nginx/conf.d/*.conf;
}
`

const nginxDefaultConfTemplate = `upstream fastcgi_backend {
    server php:9000;
}

server {
    listen 80;
    server_name {{.Name}}.test;

    set $MAGE_ROOT /var/www/html;
    set $MAGE_MODE developer;

    root $MAGE_ROOT/pub;

    index index.php;
    autoindex off;
    charset UTF-8;

    location /setup {
        root $MAGE_ROOT;
        location ~ ^/setup/index.php {
            fastcgi_pass fastcgi_backend;
            fastcgi_index index.php;
            fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
            include fastcgi_params;
        }

        location ~ ^/setup/(?!pub/). {
            deny all;
        }

        location ~ ^/setup/pub/ {
            add_header X-Frame-Options "SAMEORIGIN";
        }
    }

    location /update {
        root $MAGE_ROOT;

        location ~ ^/update/index.php {
            fastcgi_split_path_info ^(/update/index.php)(/.+)$;
            fastcgi_pass fastcgi_backend;
            fastcgi_index index.php;
            fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
            fastcgi_param PATH_INFO $fastcgi_path_info;
            include fastcgi_params;
        }

        location ~ ^/update/(?!pub/). {
            deny all;
        }

        location ~ ^/update/pub/ {
            add_header X-Frame-Options "SAMEORIGIN";
        }
    }

    location / {
        try_files $uri $uri/ /index.php$is_args$args;
    }

    location /pub/ {
        location ~ ^/pub/media/(downloadable|customer|import|theme_customization/.*\.xml) {
            deny all;
        }
        alias $MAGE_ROOT/pub/;
        add_header X-Frame-Options "SAMEORIGIN";
    }

    location /static/ {
        expires max;

        location ~ ^/static/version {
            rewrite ^/static/(version\d*/)?(.*)$ /static/$2 last;
        }

        location ~* \.(ico|jpg|jpeg|png|gif|svg|js|css|swf|eot|ttf|otf|woff|woff2)$ {
            add_header Cache-Control "public";
            add_header X-Frame-Options "SAMEORIGIN";
            expires +1y;

            if (!-f $request_filename) {
                rewrite ^/static/(version\d*/)?(.*)$ /static.php?resource=$2 last;
            }
        }

        location ~* \.(zip|gz|gzip|bz2|csv|xml)$ {
            add_header Cache-Control "no-store";
            add_header X-Frame-Options "SAMEORIGIN";
            expires off;

            if (!-f $request_filename) {
               rewrite ^/static/(version\d*/)?(.*)$ /static.php?resource=$2 last;
            }
        }

        if (!-f $request_filename) {
            rewrite ^/static/(version\d*/)?(.*)$ /static.php?resource=$2 last;
        }

        add_header X-Frame-Options "SAMEORIGIN";
    }

    location /media/ {
        try_files $uri $uri/ /get.php$is_args$args;

        location ~ ^/media/theme_customization/.*\.xml {
            deny all;
        }

        location ~* \.(ico|jpg|jpeg|png|gif|svg|js|css|swf|eot|ttf|otf|woff|woff2)$ {
            add_header Cache-Control "public";
            add_header X-Frame-Options "SAMEORIGIN";
            expires +1y;
            try_files $uri $uri/ /get.php$is_args$args;
        }

        location ~* \.(zip|gz|gzip|bz2|csv|xml)$ {
            add_header Cache-Control "no-store";
            add_header X-Frame-Options "SAMEORIGIN";
            expires off;
            try_files $uri $uri/ /get.php$is_args$args;
        }

        add_header X-Frame-Options "SAMEORIGIN";
    }

    location /media/customer/ {
        deny all;
    }

    location /media/downloadable/ {
        deny all;
    }

    location /media/import/ {
        deny all;
    }

    location ~ /media/theme_customization/.*\.xml$ {
        deny all;
    }

    location ~ cron\.php {
        deny all;
    }

    location ~ (index|get|static|report|404|503|health_check)\.php$ {
        try_files $uri =404;
        fastcgi_pass fastcgi_backend;
        fastcgi_buffers 1024 4k;

        fastcgi_param PHP_FLAG "session.auto_start=off \n suhosin.session.cryptua=off";
        fastcgi_param PHP_VALUE "memory_limit=2G \n max_execution_time=18000";
        fastcgi_read_timeout 600s;
        fastcgi_connect_timeout 600s;

        fastcgi_index index.php;
        fastcgi_param SCRIPT_FILENAME $document_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ \.php$ {
        deny all;
    }
}
`

const phpDockerfileTemplate = `ARG PHP_VERSION
FROM php:$` + `{PHP_VERSION}-fpm-alpine

ARG INSTALL_OPENSWOOLE=false

RUN apk add --no-cache \
    freetype-dev \
    libjpeg-turbo-dev \
    libpng-dev \
    libxml2-dev \
    libxslt-dev \
    libzip-dev \
    icu-dev \
    oniguruma-dev \
    bash \
    git \
    patch \
    $` + `{INSTALL_OPENSWOOLE:+postgresql-dev} \
    $` + `{INSTALL_OPENSWOOLE:+autoconf} \
    $` + `{INSTALL_OPENSWOOLE:+g++} \
    $` + `{INSTALL_OPENSWOOLE:+make} \
    && docker-php-ext-configure gd --with-freetype --with-jpeg \
    && docker-php-ext-install -j$(nproc) \
        bcmath \
        gd \
        intl \
        mbstring \
        opcache \
        pdo_mysql \
        soap \
        sockets \
        xsl \
        zip

# Install OpenSwoole if enabled
RUN if [ "$` + `INSTALL_OPENSWOOLE" = "true" ]; then \
        pecl install openswoole \
        && docker-php-ext-enable openswoole \
        && echo "openswoole.use_shortname = 'Off'" >> /usr/local/etc/php/conf.d/docker-php-ext-openswoole.ini; \
    fi

# Install Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

WORKDIR /var/www/html

CMD ["php-fpm"]
`

const phpIniTemplate = `memory_limit = 4G
max_execution_time = 1800
zlib.output_compression = On
upload_max_filesize = 256M
post_max_size = 256M

opcache.enable = 1
opcache.enable_cli = 1
opcache.memory_consumption = 512
opcache.interned_strings_buffer = 16
opcache.max_accelerated_files = 100000
opcache.validate_timestamps = 1
opcache.revalidate_freq = 2
opcache.save_comments = 1
`

const phpFpmConfTemplate = `[www]
user = www-data
group = www-data
listen = 9000
pm = dynamic
pm.max_children = 50
pm.start_servers = 10
pm.min_spare_servers = 5
pm.max_spare_servers = 20
pm.max_requests = 500
`

const mariadbConfTemplate = `[mysqld]
innodb_buffer_pool_size = 1G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 2
innodb_flush_method = O_DIRECT
max_allowed_packet = 256M
table_open_cache = 4096
query_cache_type = 0
query_cache_size = 0
`

type TemplateData struct {
	config.DeckConfig
}

func GenerateDockerFiles(cfg *config.DeckConfig, deckDir string) error {
	// Create directories
	dirs := []string{
		filepath.Join(deckDir, "nginx"),
		filepath.Join(deckDir, "php"),
		filepath.Join(deckDir, "mariadb"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	data := TemplateData{DeckConfig: *cfg}

	// Generate docker-compose.yml
	if err := generateFile(filepath.Join(deckDir, "docker-compose.yml"), dockerComposeTemplate, data); err != nil {
		return err
	}

	// Generate nginx configs
	if err := generateFile(filepath.Join(deckDir, "nginx", "nginx.conf"), nginxConfTemplate, data); err != nil {
		return err
	}
	if err := generateFile(filepath.Join(deckDir, "nginx", "default.conf"), nginxDefaultConfTemplate, data); err != nil {
		return err
	}

	// Generate PHP configs
	if err := generateFile(filepath.Join(deckDir, "php", "Dockerfile"), phpDockerfileTemplate, data); err != nil {
		return err
	}
	if err := generateFile(filepath.Join(deckDir, "php", "php.ini"), phpIniTemplate, data); err != nil {
		return err
	}
	if err := generateFile(filepath.Join(deckDir, "php", "php-fpm.conf"), phpFpmConfTemplate, data); err != nil {
		return err
	}

	// Generate MariaDB config
	if err := generateFile(filepath.Join(deckDir, "mariadb", "my.cnf"), mariadbConfTemplate, data); err != nil {
		return err
	}

	return nil
}

func generateFile(path, tmplStr string, data interface{}) error {
	tmpl, err := template.New(filepath.Base(path)).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("failed to parse template for %s: %w", path, err)
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", path, err)
	}
	defer f.Close()

	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", path, err)
	}

	return nil
}
