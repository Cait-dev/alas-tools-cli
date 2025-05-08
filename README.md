# Alas-Tools-Cli

**Alas-Tools-Cli** es una herramienta de línea de comandos (CLI) para la corrección de coordenadas y la optimización de rutas de entrega. Ideal para operaciones logísticas que manejan pallets y necesitan mejorar precisión y eficiencia.

## Características

- ✅ Corrección de coordenadas **X & Y** usando **Google Places**
- 🚚 Visualización de rutas optimizadas para **pallets**
- 🗺️ Extracción de coordenadas y generación de **mapas HTML interactivos**

## Instalación

### macOS y Linux

Ejecuta el siguiente comando en tu terminal:

```bash
curl -sSL https://raw.githubusercontent.com/Cait-dev/alas-tools-cli/main/scripts/install.sh | bash
```

### Windows

Descarga el instalador desde la [página de releases](https://github.com/Cait-dev/alas-tools-cli/releases) o ejecuta en PowerShell:

```powershell
powershell -Command "iwr -useb https://raw.githubusercontent.com/Cait-dev/alas-tools-cli/main/scripts/install.ps1 | iex"
```
## Configuración

### Credenciales de API

Esta aplicación requiere credenciales para acceder a la API de Alas Express. Por razones de seguridad, estas credenciales no están incluidas en el código fuente y deben configurarse como variables de entorno.

#### Opción 1: Variables de entorno

Configura las siguientes variables de entorno en tu sistema:

```bash
# En Linux/macOS
export ALAS_API_USER="tu_usuario"
export ALAS_API_PASSWORD="tu_contraseña"

# En Windows (PowerShell)
$env:ALAS_API_USER="tu_usuario"
$env:ALAS_API_PASSWORD="tu_contraseña"

# En Windows (CMD)
set ALAS_API_USER=tu_usuario
set ALAS_API_PASSWORD=tu_contraseña
```

#### Opción 2: Archivo .env (para desarrollo)

1. Copia el archivo `.env.example` a `.env`:
   ```
   cp .env.example .env
   ```

2. Edita el archivo `.env` y añade tus credenciales reales:
   ```
   ALAS_API_USER=tu_usuario
   ALAS_API_PASSWORD=tu_contraseña
   ```

3. La aplicación cargará automáticamente estas variables desde el archivo `.env` si está presente.

> ⚠️ **IMPORTANTE**: Nunca compartas tus credenciales ni subas el archivo `.env` a GitHub u otros repositorios públicos.


## Contribución

¡Contribuciones son bienvenidas! Abre un _issue_ o un _pull request_ en el [repositorio oficial](https://github.com/Cait-dev/alas-tools-cli).

## Licencia

Este proyecto está licenciado bajo la licencia MIT. Consulta el archivo `LICENSE` para más información.
