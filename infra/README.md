# Infra — Server Configuration

Konfigurasi server untuk deployment LinkChat di production.

## Struktur

```
infra/
├── nginx/
│   └── linkchat.conf       # Nginx reverse proxy config
└── systemd/
    ├── linkchat-api.service # Go backend service
    └── linkchat-web.service # Next.js frontend service
```

## Setup di Server (Ubuntu/Debian)

### 1. Buat user linkchat

```bash
sudo useradd -r -s /bin/false -d /opt/linkchat linkchat
sudo mkdir -p /opt/linkchat
sudo chown linkchat:linkchat /opt/linkchat
```

### 2. Deploy aplikasi

```bash
# Copy project ke server
sudo cp -r /path/to/linkchat /opt/linkchat

# Build Go backend
cd /opt/linkchat/services/api
sudo -u linkchat go build -o bin/api ./cmd/main.go

# Build Next.js
cd /opt/linkchat/apps/web
sudo -u linkchat npm ci
sudo -u linkchat npm run build
```

### 3. Konfigurasi environment

```bash
# Copy dan edit .env
sudo cp /opt/linkchat/.env.example /opt/linkchat/.env
sudo nano /opt/linkchat/.env
# Isi DB_PASSWORD, SESSION_SECRET, dan lainnya
```

### 4. Install systemd services

```bash
sudo cp infra/systemd/linkchat-api.service /etc/systemd/system/
sudo cp infra/systemd/linkchat-web.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable linkchat-api linkchat-web
sudo systemctl start linkchat-api linkchat-web
```

### 5. Setup Nginx

```bash
# Install nginx
sudo apt install nginx

# Copy config
sudo cp infra/nginx/linkchat.conf /etc/nginx/sites-available/linkchat
sudo ln -s /etc/nginx/sites-available/linkchat /etc/nginx/sites-enabled/
sudo nginx -t
sudo nginx -s reload
```

### 6. SSL dengan Certbot

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d linkchat.yourdomain.com
```

## Commands Berguna

```bash
# Status services
sudo systemctl status linkchat-api linkchat-web

# Lihat logs
sudo journalctl -u linkchat-api -f
sudo journalctl -u linkchat-web -f

# Restart
sudo systemctl restart linkchat-api linkchat-web

# Reload nginx
sudo nginx -s reload
```
