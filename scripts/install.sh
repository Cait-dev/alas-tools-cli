#!/bin/bash

# Colores para outputs
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Detectar sistema operativo
OS="$(uname -s)"
ARCH="$(uname -m)"

echo -e "${GREEN}Instalando Alas-Tools-Cli...${NC}"

# Carpeta de destino
if [ "$OS" = "Darwin" ]; then
    # macOS
    INSTALL_DIR="$HOME/Applications/AlasCli"
    if [ "$ARCH" = "arm64" ]; then
        BINARY_URL="https://github.com/cait-dev/alas-tools-cli/releases/latest/download/alas-tools-cli-mac-arm64"
        BINARY_NAME="alas-tools-cli-mac-arm64"
    else
        BINARY_URL="https://github.com/cait-dev/alas-tools-cli/releases/latest/download/alas-tools-cli-mac"
        BINARY_NAME="alas-tools-cli-mac"
    fi
elif [ "$OS" = "Linux" ]; then
    # Linux
    INSTALL_DIR="$HOME/.local/bin"
    BINARY_URL="https://github.com/cait-dev/alas-tools-cli/releases/latest/download/alas-tools-cli-linux"
    BINARY_NAME="alas-tools-cli-linux"
else
    echo -e "${RED}Sistema operativo no soportado.${NC}"
    exit 1
fi

# Crear directorio de instalación si no existe
mkdir -p "$INSTALL_DIR"

# Descargar el binario
echo "Descargando Alas-Tools-Cli para $OS ($ARCH)..."
if command -v curl &>/dev/null; then
    curl -L "$BINARY_URL" -o "$INSTALL_DIR/$BINARY_NAME"
elif command -v wget &>/dev/null; then
    wget "$BINARY_URL" -O "$INSTALL_DIR/$BINARY_NAME"
else
    echo -e "${RED}Se requiere curl o wget para la instalación.${NC}"
    exit 1
fi

# Verificar si la descarga fue exitosa
if [ $? -ne 0 ]; then
    echo -e "${RED}Error al descargar el binario.${NC}"
    exit 1
fi

# Hacer ejecutable el binario
chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Crear un symlink en un directorio del PATH (opcional)
if [ "$OS" = "Darwin" ]; then
    # macOS
    if [ ! -d "$HOME/bin" ]; then
        mkdir -p "$HOME/bin"
    fi
    ln -sf "$INSTALL_DIR/$BINARY_NAME" "$HOME/bin/alas-cli"
    
    # Verificar si $HOME/bin está en PATH, si no, sugerir añadirlo
    if [[ ":$PATH:" != *":$HOME/bin:"* ]]; then
        echo "Añade la siguiente línea a tu archivo .bash_profile o .zshrc:"
        echo 'export PATH="$HOME/bin:$PATH"'
    fi
elif [ "$OS" = "Linux" ]; then
    # Linux
    ln -sf "$INSTALL_DIR/$BINARY_NAME" "$INSTALL_DIR/alas-cli"
    
    # Verificar si $HOME/.local/bin está en PATH
    if [[ ":$PATH:" != *":$HOME/.local/bin:"* ]]; then
        echo "Añade la siguiente línea a tu archivo .bashrc o .profile:"
        echo 'export PATH="$HOME/.local/bin:$PATH"'
    fi
fi

echo -e "${GREEN}¡Instalación completada!${NC}"
echo "Puedes ejecutar Alas-Tools-Cli con el comando 'alas-cli'"