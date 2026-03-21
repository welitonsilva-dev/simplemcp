# =============================================================
#  cli.ps1 — Wrapper Windows para o MCP CLI
#  Detecta disco e diretório atual automaticamente
#  Suporta múltiplos discos (C:\, D:\, E:\, etc.)
# =============================================================

param(
    [Parameter(ValueFromRemainingArguments=$true)]
    [string[]]$Args
)

# Detecta o disco e path atual
$currentLocation = Get-Location
$currentDrive    = $currentLocation.Drive.Root        # Ex: C:\  ou  D:\
$currentPath     = $currentLocation.Path              # Ex: C:\Users\dev\projetos\repo1
$driveLetter     = $currentLocation.Drive.Name.ToLower() # Ex: c, d, e

# No container, cada disco é montado em /app/host/<letra>
# Ex: C:\ → /app/host/c    D:\ → /app/host/d
$containerDriveRoot = "/app/host/$driveLetter"

# Converte o path atual para o equivalente no container
# Ex: C:\Users\dev\projetos → /app/host/c/Users/dev/projetos
$relativePath = $currentPath.Substring($currentDrive.Length).Replace('\', '/')
$containerCWD = "$containerDriveRoot/$relativePath".TrimEnd('/')

# Detecta todos os discos disponíveis no sistema
$drives = Get-PSDrive -PSProvider FileSystem | Where-Object { $_.Root -match '^[A-Z]:\\$' }

Write-Host "🖥️  Disco atual: $currentDrive" -ForegroundColor Cyan
Write-Host "📁  Diretório:   $currentPath" -ForegroundColor Cyan
Write-Host "🐳  Container:   $containerCWD" -ForegroundColor Cyan
Write-Host ""

# Monta os argumentos de volume para cada disco encontrado
$volumeArgs = @()
foreach ($drive in $drives) {
    $letter = $drive.Name.ToLower()
    $root   = $drive.Root
    $volumeArgs += "-v"
    $volumeArgs += "${root}:/app/host/${letter}"
    Write-Host "   Mapeando $root → /app/host/$letter" -ForegroundColor DarkGray
}

Write-Host ""

# Executa o container com todos os discos mapeados
docker compose run --rm `
    -e HOST_ROOT=$currentDrive `
    -e HOST_CWD=$currentPath `
    -e CONTAINER_CWD=$containerCWD `
    -e HOST_OS=windows `
    -w $containerCWD `
    @volumeArgs `
    humancli-server @Args