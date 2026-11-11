# ─── LinkChat Development Startup Script (Windows PowerShell) ────────────────
# Usage: .\scripts\dev.ps1
# Starts both backend and frontend in separate windows

$root = Split-Path $PSScriptRoot -Parent

Write-Host "Starting LinkChat Development Environment..." -ForegroundColor Cyan
Write-Host ""

# ── Add common tool paths ─────────────────────────────────────────────────────
$extraPaths = @(
    "C:\Program Files\Go\bin",
    "C:\Go\bin",
    "C:\Program Files\nodejs",
    "C:\Users\$env:USERNAME\AppData\Roaming\npm"
)
foreach ($p in $extraPaths) {
    if (Test-Path $p) { $env:PATH = "$p;$env:PATH" }
}

# ── Check prerequisites ───────────────────────────────────────────────────────
$goOk   = Get-Command go   -ErrorAction SilentlyContinue
$nodeOk = Get-Command node -ErrorAction SilentlyContinue

if (-not $goOk) {
    Write-Host "ERROR: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Download from: https://go.dev/dl/" -ForegroundColor Yellow
    exit 1
}

if (-not $nodeOk) {
    Write-Host "ERROR: Node.js is not installed or not in PATH" -ForegroundColor Red
    Write-Host "Download from: https://nodejs.org/" -ForegroundColor Yellow
    exit 1
}

Write-Host "Go version: $(go version)" -ForegroundColor Green
Write-Host "Node version: $(node --version)" -ForegroundColor Green
Write-Host ""

# ── Load .env ────────────────────────────────────────────────────────────────
$envFile = "$root\.env"
if (Test-Path $envFile) {
    Write-Host "Loading .env..." -ForegroundColor Gray
    Get-Content $envFile | ForEach-Object {
        if ($_ -match "^([^#][^=]*)=(.*)$") {
            [System.Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim(), "Process")
        }
    }
} else {
    Write-Host "WARNING: .env file not found. Copy .env.example to .env first." -ForegroundColor Yellow
    exit 1
}

# ── Start Go backend in new window ────────────────────────────────────────────
Write-Host "Starting Go backend on port $env:API_PORT..." -ForegroundColor Blue

$envVars = "`$env:DB_HOST='$env:DB_HOST'; `$env:DB_PORT='$env:DB_PORT'; `$env:DB_USER='$env:DB_USER'; `$env:DB_PASSWORD='$env:DB_PASSWORD'; `$env:DB_NAME='$env:DB_NAME'; `$env:API_PORT='$env:API_PORT'; `$env:SESSION_SECRET='$env:SESSION_SECRET'"

Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "`$env:PATH = 'C:\Program Files\Go\bin;' + `$env:PATH; $envVars; Set-Location '$root\services\api'; Write-Host 'LinkChat API - Starting...' -ForegroundColor Cyan; go run ./cmd/main.go"
)

# Wait a moment for backend to start
Start-Sleep -Seconds 2

# ── Start Next.js frontend in new window ──────────────────────────────────────
Write-Host "Starting Next.js frontend on port 3000..." -ForegroundColor Blue
Start-Process powershell -ArgumentList @(
    "-NoExit",
    "-Command",
    "`$env:PATH = 'C:\Program Files\nodejs;' + `$env:PATH; Set-Location '$root\apps\web'; Write-Host 'LinkChat Web' -ForegroundColor Cyan; npm run dev"
)

Write-Host ""
Write-Host "LinkChat is starting up!" -ForegroundColor Green
Write-Host "Frontend: http://localhost:3000" -ForegroundColor Cyan
Write-Host "Backend:  http://localhost:8080" -ForegroundColor Cyan
Write-Host "Health:   http://localhost:8080/health" -ForegroundColor Cyan
