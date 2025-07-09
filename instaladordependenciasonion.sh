#!/bin/bash
# Script para instalar todas las dependencias de buscaonion
# by SaturniCipher

set -e

# Verifica si python3 está instalado
if ! command -v python3 >/dev/null 2>&1; then
    echo "[!] Instalando python3..."
    sudo apt-get update && sudo apt-get install -y python3
fi

# Verifica si pip está instalado
if ! command -v pip3 >/dev/null 2>&1; then
    echo "[!] Instalando python3-pip..."
    sudo apt-get update && sudo apt-get install -y python3-pip
fi

# Instala dependencias de Python
pip3 install --upgrade pip
pip3 install requests[socks] beautifulsoup4 rich PySocks

echo "[+] Dependencias instaladas. Ya puedes usar buscaonion."
