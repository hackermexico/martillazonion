# 🧨 Martillazonion - Carga Distribuida contra Servicios .onion

`Martillazonion` es una herramienta ofensiva para pruebas de carga o denegación de servicio (DoS) sobre servicios `.onion` en la red Tor. Usa múltiples instancias de Tor sobre distintos puertos SOCKS5 para rotar conexiones y maximizar el impacto sin saturar una sola salida.

## ⚔️ Versiones incluidas

- `martillazonion.py`: Versión sencilla que permite usar un solo puerto Tor SOCKS5 (`-p`).
- `martillazonion_multi.py`: Versión mejorada que balancea ataques sobre múltiples instancias Tor.
- `tor-multi-launcher.sh`: Script para iniciar instancias de Tor en puertos 9055 al 9060.
- `tor-multi-stop.sh`: Script para detener todas esas instancias de forma limpia.

---

## 🚀 Requisitos

- Python 3.x
- Módulo `PySocks`:  
  ```bash
  pip3 install pysocks
