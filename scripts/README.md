# Scripts — Helper Scripts for Development & Deployment

## Available Scripts

| Script | Usage | Deskripsi |
|--------|-------|-----------|
| `dev.ps1` | `.\scripts\dev.ps1` | Jalankan backend + frontend sekaligus |
| `migrate.ps1` | `.\scripts\migrate.ps1 up` | Jalankan semua database migrations |
| `migrate.ps1` | `.\scripts\migrate.ps1 down` | Rollback migration terakhir |
| `migrate.ps1` | `.\scripts\migrate.ps1 status` | Lihat daftar migration files |
| `build.ps1` | `.\scripts\build.ps1` | Build production binary |
| `seed.ps1` | `.\scripts\seed.ps1` | Insert data awal (owner account) |

## Quick Start (Windows)

```powershell
# 1. Copy environment config
Copy-Item .env.example .env
# Edit .env sesuai kebutuhan

# 2. Run migrations
.\scripts\migrate.ps1 up

# 3. Seed database
.\scripts\seed.ps1

# 4. Start development
.\scripts\dev.ps1
```

## Notes

- Semua script menggunakan PowerShell (.ps1)
- Script membaca konfigurasi dari `.env` di root project
- MySQL harus sudah running (XAMPP) sebelum migrate/seed
- Go dan Node.js harus sudah terinstall sebelum dev/build
