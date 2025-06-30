#!/bin/bash

# Rango de puertos que se usaron
INICIO=9055
FIN=9060

echo "[!] Deteniendo instancias Tor en puertos $INICIO a $FIN..."

for PORT in $(seq $INICIO $FIN); do
    CONF="/etc/tor/torrc$PORT"
    LOG="/var/log/tor$PORT.log"
    DIR="/var/lib/tor$PORT"

    echo " - Terminando procesos en puerto $PORT..."
    sudo pkill -f "tor.*$CONF"

    # Opcional: eliminar archivos (descomentar si quieres borrar todo)
    # echo "   - Borrando configuración y datos..."
    # sudo rm -f "$CONF" "$LOG"
    # sudo rm -rf "$DIR"
done

echo "[✓] Todas las instancias Tor han sido detenidas."
