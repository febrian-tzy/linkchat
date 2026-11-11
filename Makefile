.PHONY: dev api-dev web-dev build api-build web-build migrate-up migrate-down seed

# ── Development ───────────────────────────────────────────────
dev:
	powershell -ExecutionPolicy Bypass -File scripts/dev.ps1

api-dev:
	cd services/api && go run ./cmd/main.go

web-dev:
	cd apps/web && npm run dev

# ── Build ─────────────────────────────────────────────────────
build:
	powershell -ExecutionPolicy Bypass -File scripts/build.ps1

api-build:
	cd services/api && go build -o bin/api ./cmd/main.go

web-build:
	cd apps/web && npm run build

# ── Database ──────────────────────────────────────────────────
migrate-up:
	powershell -ExecutionPolicy Bypass -File scripts/migrate.ps1 up

migrate-down:
	powershell -ExecutionPolicy Bypass -File scripts/migrate.ps1 down

migrate-status:
	powershell -ExecutionPolicy Bypass -File scripts/migrate.ps1 status

seed:
	powershell -ExecutionPolicy Bypass -File scripts/seed.ps1
