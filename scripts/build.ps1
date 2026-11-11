# ─── LinkChat Build Script (Windows PowerShell) ──────────────────────────────
# Builds both backend and frontend for production deployment
# Usage: .\scripts\build.ps1

$root = Split-Path $PSScriptRoot -Parent

Write-Host "Building LinkChat for production..." -ForegroundColor Cyan
Write-Host ""

# ── Build Go backend ─────────────────────────────────────────────────────────
Write-Host "[1/2] Building Go backend..." -ForegroundColor Blue
$apiDir = "$root\services\api"
$outDir = "$apiDir\bin"

if (-not (Test-Path $outDir)) {
    New-Item -ItemType Directory -Path $outDir | Out-Null
}

Set-Location $apiDir
$env:GOOS   = "linux"
$env:GOARCH = "amd64"
go build -ldflags="-s -w" -o "$outDir\api" ./cmd/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "Go build FAILED" -ForegroundColor Red
    exit 1
}

Write-Host "Go backend built: $outDir\api" -ForegroundColor Green
Write-Host ""

# ── Build Next.js frontend ────────────────────────────────────────────────────
Write-Host "[2/2] Building Next.js frontend..." -ForegroundColor Blue
$webDir = "$root\apps\web"
Set-Location $webDir

npm run build

if ($LASTEXITCODE -ne 0) {
    Write-Host "Next.js build FAILED" -ForegroundColor Red
    exit 1
}

Write-Host "Next.js frontend built: $webDir\.next" -ForegroundColor Green
Write-Host ""
Write-Host "Build complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Deploy files:"
Write-Host "  Backend:  services/api/bin/api"
Write-Host "  Frontend: apps/web/.next + apps/web/package.json"
