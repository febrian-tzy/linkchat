# ─── LinkChat Database Migration Script (Windows PowerShell) ─────────────────
# Usage:
#   .\scripts\migrate.ps1 up      # Run all pending migrations
#   .\scripts\migrate.ps1 down    # Rollback last migration
#   .\scripts\migrate.ps1 status  # Show migration status

param(
    [Parameter(Position=0)]
    [ValidateSet("up", "down", "status")]
    [string]$Action = "up"
)

$root     = Split-Path $PSScriptRoot -Parent
$mysql    = "C:\xampp\mysql\bin\mysql.exe"
$migrDir  = "$root\database\migrations"

# ── Load .env ────────────────────────────────────────────────────────────────
$envFile = "$root\.env"
if (-not (Test-Path $envFile)) {
    Write-Host "ERROR: .env not found. Copy .env.example to .env first." -ForegroundColor Red
    exit 1
}

Get-Content $envFile | ForEach-Object {
    if ($_ -match "^([^#][^=]*)=(.*)$") {
        [System.Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim(), "Process")
    }
}

$dbUser = $env:DB_USER
$dbPass = $env:DB_PASSWORD
$dbName = $env:DB_NAME
$dbHost = $env:DB_HOST
$dbPort = $env:DB_PORT

if (-not (Test-Path $mysql)) {
    Write-Host "ERROR: mysql.exe not found at $mysql" -ForegroundColor Red
    exit 1
}

function Run-SQL {
    param([string]$file)
    Write-Host "  Running: $(Split-Path $file -Leaf)" -ForegroundColor Gray
    $args = @("-u", $dbUser, "-h", $dbHost, "--port=$dbPort", $dbName)
    if ($dbPass) { $args = @("-u", $dbUser, "-p$dbPass", "-h", $dbHost, "--port=$dbPort", $dbName) }
    Get-Content $file | & $mysql @args 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Host "  FAILED: $(Split-Path $file -Leaf)" -ForegroundColor Red
        return $false
    }
    Write-Host "  OK: $(Split-Path $file -Leaf)" -ForegroundColor Green
    return $true
}

switch ($Action) {
    "up" {
        Write-Host "Running migrations UP..." -ForegroundColor Cyan
        $files = Get-ChildItem "$migrDir\*.up.sql" | Sort-Object Name
        foreach ($f in $files) {
            Run-SQL $f.FullName
        }
        Write-Host "Done." -ForegroundColor Green
    }

    "down" {
        Write-Host "Running last migration DOWN..." -ForegroundColor Yellow
        $files = Get-ChildItem "$migrDir\*.down.sql" | Sort-Object Name -Descending
        if ($files.Count -gt 0) {
            Run-SQL $files[0].FullName
        }
        Write-Host "Done." -ForegroundColor Green
    }

    "status" {
        Write-Host "Migration files in $migrDir :" -ForegroundColor Cyan
        Get-ChildItem "$migrDir\*.up.sql" | Sort-Object Name | ForEach-Object {
            Write-Host "  $($_.Name)" -ForegroundColor Gray
        }
    }
}
