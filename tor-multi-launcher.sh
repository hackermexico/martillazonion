#!/bin/bash

# Rango de puertos
INICIO=9055
FIN=9060

echo "[+] Iniciando instancias Tor en puertos $INICIO a $FIN..."

for PORT in $(seq $INICIO $FIN); do
    DIR="/var/lib/tor$PORT"
    LOG="/var/log/tor$PORT.log"
    CONF="/etc/tor/torrc$PORT"

    # Crear directorio de datos si no existe
    if [ ! -d "$DIR" ]; then
        sudo mkdir -p "$DIR"
        sudo chown debian-tor:debian-tor "$DIR"
        sudo chmod 700 "$DIR"
    fi

    # Crear archivo de configuración torrc específico
    sudo bash -c "cat > $CONF" <<EOF
SocksPort 127.0.0.1:$PORT
DataDirectory $DIR
Log notice file $LOG
EOF

    # Lanzar instancia de Tor con esa configuración
    echo " - Iniciando Tor en puerto $PORT..."
    sudo -u debian-tor tor -f "$CONF" &
    sleep 1
done

echo "[✓] Todas las instancias de Tor han sido lanzadas."
