# Script de instalación para Alas-Tools-Cli en Windows

Write-Host "Instalando Alas-Tools-Cli para Windows..." -ForegroundColor Green

# Crear directorio de instalación
$installDir = "$env:USERPROFILE\AlasCli"
if (-not (Test-Path -Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

# Detectar arquitectura
$arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }

# URL de descarga
$binaryUrl = if ($arch -eq "386") {
    "https://github.com/cait-dev/alas-tools-cli/releases/latest/download/alas-tools-cli-386.exe"
} else {
    "https://github.com/cait-dev/alas-tools-cli/releases/latest/download/alas-tools-cli.exe"
}
$binaryName = if ($arch -eq "386") { "alas-tools-cli-386.exe" } else { "alas-tools-cli.exe" }

# Descargar el binario
Write-Host "Descargando desde $binaryUrl..."
try {
    [Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
    Invoke-WebRequest -Uri $binaryUrl -OutFile "$installDir\$binaryName"
} catch {
    Write-Host "Error al descargar el binario: $_" -ForegroundColor Red
    exit 1
}

# Crear un acceso directo en el Escritorio
Write-Host "Creando acceso directo..."
$WshShell = New-Object -ComObject WScript.Shell
$Shortcut = $WshShell.CreateShortcut("$env:USERPROFILE\Desktop\Alas-Tools-Cli.lnk")
$Shortcut.TargetPath = "$installDir\$binaryName"
$Shortcut.Save()

# Preguntar si quiere añadir al PATH
$addToPath = Read-Host "¿Desea añadir Alas-Tools-Cli al PATH del sistema? (S/N)"
if ($addToPath -eq "S" -or $addToPath -eq "s") {
    Write-Host "Añadiendo al PATH..."
    $currentPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if (-not $currentPath.Contains($installDir)) {
        [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$installDir", "User")
        Write-Host "Será necesario reiniciar PowerShell para usar el comando." -ForegroundColor Yellow
    }
}

Write-Host "`nInstalación completada." -ForegroundColor Green
Write-Host "Puede ejecutar Alas-Tools-Cli desde el acceso directo en el Escritorio."
if ($addToPath -eq "S" -or $addToPath -eq "s") {
    Write-Host "O usando '$binaryName' en una nueva consola de PowerShell."
}

Write-Host "`nPresione Enter para salir..."
Read-Host