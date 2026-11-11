# LinkChat

Platform komunikasi pelanggan berbasis SaaS.

## Struktur Project

```
linkchat/
├── apps/web          # Frontend Next.js (TypeScript + Tailwind)
├── services/api      # Backend Go
├── database/
│   └── migrations/   # SQL migrations MySQL
├── infra/            # Nginx, systemd configs
├── scripts/          # Helper scripts
└── docs/             # Dokumentasi
```

## Prasyarat

- Node.js >= 18
- Go >= 1.22
- MySQL >= 8.0 (atau XAMPP)

## Setup

### 1. Clone & konfigurasi environment

```bash
cp .env.example .env
# Edit .env sesuai environment lokal
```

### 2. Jalankan migrasi database

Buka phpMyAdmin atau mysql CLI, buat database `linkchat`, lalu jalankan:

```sql
source database/migrations/000001_create_tenants.up.sql
source database/migrations/000002_create_tenant_users.up.sql
```

Atau gunakan Makefile:

```bash
make migrate-up
```

### 3. Backend

```bash
make api-dev
# atau manual:
cd services/api && go run ./cmd/main.go
```

### 4. Frontend

```bash
make web-dev
# atau manual:
cd apps/web && npm install && npm run dev
```

## API Endpoints (Tahap Pertama)

| Method | Endpoint                  | Deskripsi                  |
|--------|---------------------------|----------------------------|
| POST   | /api/v1/auth/register     | Registrasi tenant + admin  |
| POST   | /api/v1/auth/login        | Login tenant               |
| POST   | /api/v1/auth/logout       | Logout (hapus session)     |
| GET    | /api/v1/auth/me           | Get user & tenant aktif    |

## Frontend Routes

| Route          | Deskripsi                        |
|----------------|----------------------------------|
| /              | Home (pilihan login/register)    |
| /register      | Form registrasi tenant baru      |
| /login/tenant  | Login tenant                     |
| /login/owner   | Login Owner (placeholder)        |
| /app           | Dashboard (protected)            |
