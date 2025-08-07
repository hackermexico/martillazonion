# 🛠️ Martillazonion - Carga Distribuida contra Servicios .onion

`Martillazonion` es una herramienta ofensiva para pruebas de carga o denegación de servicio (DoS) sobre servicios `.onion` en la red Tor. Usa múltiples instancias de Tor sobre distintos puertos SOCKS5 para rotar conexiones y maximizar el impacto sin saturar una sola salida.

## ⚔️ Versiones incluidas, la versión ALPHA es en la que se está experimentando nuevas técnicas.

- `martillazonion.py`: Versión sencilla que permite usar un solo puerto Tor SOCKS5 (`-p`).
- `martillazonion_multi.py`: Versión mejorada que balancea ataques sobre múltiples instancias Tor.
- `tor-multi-launcher.sh`: Script para iniciar instancias de Tor en puertos 9055 al 9060.
- `tor-multi-stop.sh`: Script para detener todas esas instancias de forma limpia.
- `browser_onion.py`: Script para buscar sitios .onion desde la terminal.
- `buscadonion.py`: Script para metabuscar en casi todos los buscadores de la deepweb.
- `crawleronion.py`: Script para rastrear sitios accesibles únicamente a través de la red Tor. Extrae los enlaces presentes en una página .onion, guarda los resultados en archivos .json y .csv, y presenta la información de forma organizada y automatizada.
---

## 🚀 Requisitos

- Python 3.x
- Módulo `PySocks`:  
  ```bash
  pip3 install pysocks

  # 💣 Martillazonion ⚔️
Ataques de carga distribuidos contra servicios `.onion` usando múltiples instancias Tor SOCKS5.

---

## 🚀 ¿Qué es?

**Martillazonion** es una herramienta ofensiva diseñada para realizar ataques de denegación de servicio (DoS) sobre servicios `.onion` en la red Tor. Usa la rotación de instancias Tor en puertos distintos para evitar cuellos de botella y superar las limitaciones del protocolo Tor al realizar múltiples conexiones simultáneas.

---

## 🧰 ¿Por qué múltiples instancias de Tor?

Tor tiene **límites internos de circuitos y conexiones por instancia**. Si haces demasiadas peticiones simultáneas (más de ~1000 hilos), el proceso `tor` puede:

- Saturarse y cerrar sockets.
- Lanzar errores de "assertion failed".
- Tirar conexiones sin avisar.

La solución: **crear múltiples instancias de Tor**, cada una en un puerto SOCKS5 distinto. Esto permite:

- Balancear la carga de hilos entre procesos.
- Aumentar la presión del ataque.
- Evitar que el Tor principal (9050) se caiga.

---

## 🧱 Estructura del proyecto

martillazonion/
├── martillazonion.py # Versión sencilla (un solo puerto SOCKS)
├── martillazonion_multi.py # Versión multipuerto (9055–9060)
├── tor-multi-launcher.sh # Crea instancias de Tor en puertos 9055-9060
├── tor-multi-stop.sh # Detiene y limpia las instancias Tor
└── README.md

COMO USAR:

python3 martillazonion.py <dominio.onion> -p 9050 -t 600 -r 150

-p: Puerto SOCKS5 (por defecto 9050)

-t: Número máximo de hilos (por defecto 500)

-r: Hilos por ráfaga (por defecto 100)

🛠️ Tips para tunear
En la versión sencilla:

Usa -t hasta 800 hilos si solo usarás un tor en 9050.

Si te truena, reduce a 600 y usa más nodos.

En la versión multipuerto:

Aumenta el rango de puertos en los scripts (ej. 9055–9070).

Aumenta RAFA_MAX o PAQUETES_POR_CONN en el .py multipuerto.

Usa htop o ps aux | grep tor para monitorear carga por proceso.

✅ Opción 2: Versión multipuerto (multi instancia de Tor)
1. Iniciar instancias de Tor en puertos 9055–9060:

sudo ./tor-multi-launcher.sh
Esto crea múltiples torrc, DataDirectory y lanza 6 procesos tor independientes.

2. Ejecutar ataque:

python3 martillazonion_multi.py
Se te pedirá el dominio .onion y el script comenzará a rotar automáticamente entre los puertos 9055 a 9060 para lanzar hilos distribuidos.

⚠️ Advertencia
Este proyecto es solo para:

Auditorías controladas en servicios propios.

Simulaciones de resistencia en laboratorios.

# buscaonion

Buscador y fuzzer avanzado de sitios .onion (deep web) con soporte multimotor, escaneo de puertos, dorks y detección de tecnologías. Ideal para OSINT, pentesting y exploración h4x0r.

## Características
- Búsqueda en más de 10 motores onion (ahmia, darksearch, phobos, onionland, tordex, haystak, torch, notevil, kilos, recon, candle, deepsearch, onionsearch, darknetlive).
- Fuzzing/dorks automáticos sobre cada .onion encontrado (admin, login, backup, .git, etc).
- Escaneo de puertos básicos (tipo nmap) y banner grabber.
- Detección de tecnologías (Apache, nginx, PHP, etc).
- Feedback en tiempo real y resumen de hallazgos.
- Exporta resultados a `resultados_onion_fuzz.txt`.
- Ctrl+C para cerrar y guardar todo.

## Uso rápido
1. Instala dependencias:
   ```bash
   bash instalar_dependencias_onion.sh
   ```
2. Ejecuta el script:
   ```bash
   python3 browser_onion.py
   ```
3. Ingresa una o varias palabras clave (separadas por coma) y deja que la tool haga el resto.

## Requisitos
- Python 3
- Tor corriendo en 127.0.0.1:9050
- Linux recomendado

## Autor
by SaturniCipher basado en codigo debugsec

---

# browser_onion.py (versión básica)

Script simple para buscar sitios .onion usando Ahmia y mostrar resultados en consola.

## Uso
1. Instala dependencias:
   ```bash
   pip3 install requests[socks] beautifulsoup4 rich
   ```
2. Ejecuta:
   ```bash
   python3 browser_onion.py
   ```
3. Ingresa una palabra clave y verás los .onion encontrados.

## Características
- Busca en Ahmia por keyword.
- Muestra resultados en consola con colores.
- Exporta a `resultados_onion_fuzz.txt`.
- Ctrl+C para guardar y salir.

## Autor
by DebugSec

Experimentación en entornos legales y educativos.

NO lo uses para atacar servicios de terceros. El mal uso de esta herramienta es tu total responsabilidad.
