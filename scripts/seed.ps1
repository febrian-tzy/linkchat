# ─── LinkChat Database Seed Script (Windows PowerShell) ──────────────────────
# Creates initial owner account and sample data for development
# Usage: .\scripts\seed.ps1

$root  = Split-Path $PSScriptRoot -Parent
$mysql = "C:\xampp\mysql\bin\mysql.exe"

# ── Load .env ────────────────────────────────────────────────────────────────
Get-Content "$root\.env" | ForEach-Object {
    if ($_ -match "^([^#][^=]*)=(.*)$") {
        [System.Environment]::SetEnvironmentVariable($matches[1].Trim(), $matches[2].Trim(), "Process")
    }
}

$dbUser = $env:DB_USER
$dbName = $env:DB_NAME
$dbHost = $env:DB_HOST

Write-Host "Seeding database '$dbName'..." -ForegroundColor Cyan

# Run seed SQL
$sql = @"
-- Seed: Default Owner Account
-- Password: owner123
INSERT IGNORE INTO owners (full_name, email, username, password_hash)
VALUES (
    'Platform Owner',
    'owner@linkchat.id',
    'owner',
    '\$2a\$12\$u4Vy9CgYJ4pmtW3GZb.Y4uZhIC0YoGGNXdDmLW/VWiJus2f0YJn7u'
);
SELECT 'Owner account ready: owner / owner123' AS info;
"@

$sql | & $mysql -u $dbUser -h $dbHost $dbName 2>&1

Write-Host "Seed complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Default accounts:"
Write-Host "  Owner:  username=owner  password=owner123"
Write-Host "  Tenant: Register at http://localhost:3000/register"
