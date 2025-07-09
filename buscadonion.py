#!/usr/bin/env python3
import subprocess
import sys
import os
import time
import random
import re
import socket
import socks

# === AUTOINSTALACIÓN DE DEPENDENCIAS ===
try:
    import requests
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", "requests[socks]"])
    import requests

try:
    from bs4 import BeautifulSoup
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", "beautifulsoup4"])
    from bs4 import BeautifulSoup

try:
    from rich.console import Console
    from rich.table import Table
    from rich.panel import Panel
    from rich.text import Text
except ImportError:
    subprocess.check_call([sys.executable, "-m", "pip", "install", "rich"])
    from rich.console import Console
    from rich.table import Table
    from rich.panel import Panel
    from rich.text import Text

console = Console()

# === CONFIGURACIÓN DE TOR Y STEALTH ===
TOR_SOCKS = os.environ.get("TOR_SOCKS", "127.0.0.1:9050")
USER_AGENTS = [
    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
    "curl/7.87.0", "Wget/1.21.1", "Lynx/2.8.9rel.1",
    "python-requests/2.31", "Go-http-client/1.1", "libwww-perl/6.66"
]

# === LISTA DE QUERIES COMUNES ===
QUERIES_COMUNES = [
    "market", "forum", "drugs", "carding", "bitcoin", "hacking", "escrow", "mail", "hosting", "search", "wiki", "chat", "vpn", "exploit", "leak", "paste", "crypto", "shop", "guns", "porn", "fraud", "bank", "cc", "wallet", "exchange"
]

# === VALIDACIÓN DE DOMINIOS .ONION V3 ===
ONION_V3_REGEX = re.compile(r"[a-z2-7]{56}\.onion", re.IGNORECASE)
ONION_V2_REGEX = re.compile(r"[a-z2-7]{16}\.onion", re.IGNORECASE)

def es_onion_v3(url):
    return bool(ONION_V3_REGEX.search(url))

def es_onion_v2(url):
    return bool(ONION_V2_REGEX.search(url))

# === BANNER ===
def mostrar_banner():
    console.print("\n[bold cyan]buscaonion — by SaturniCipher[/bold cyan]\n", style="bold")

# === BÚSQUEDA Y SCRAPING POR TOR ===
def buscar_onion(termino, usar_tor=True, verbose=True):
    if verbose:
        console.print(f"[bold yellow]\n🔍 Buscando sitios .onion relacionados con:[/bold yellow] [bold cyan]{termino}[/bold cyan]")
    url = f"https://ahmia.fi/search/?q={termino}"
    headers = {
        "User-Agent": random.choice(USER_AGENTS)
    }
    proxies = {}
    if usar_tor:
        host, port = TOR_SOCKS.split(":")
        proxies = {
            "http": f"socks5h://{host}:{port}",
            "https": f"socks5h://{host}:{port}"
        }
    try:
        respuesta = requests.get(url, headers=headers, timeout=20, proxies=proxies)
        respuesta.raise_for_status()
    except Exception as e:
        print(f"[!] Error al conectar con Ahmia: {e}")
        return []
    soup = BeautifulSoup(respuesta.text, "html.parser")
    enlaces = set()
    v3_count, v2_count, otros = 0, 0, 0
    for a in soup.find_all("a", href=True):
        href = a["href"].strip()
        if ".onion" in href:
            if es_onion_v3(href):
                v3_count += 1
                if verbose:
                    console.print(f"[bold green][v3][+] Sitio v3 encontrado:[/bold green] [white]{href}[/white]")
                enlaces.add(href)
            elif es_onion_v2(href):
                v2_count += 1
                if verbose:
                    console.print(f"[bold yellow][v2][!] Dominio v2 obsoleto:[/bold yellow] [white]{href}[/white]")
            else:
                otros += 1
                if verbose:
                    console.print(f"[bold red][?][!] Posible falso positivo:[/bold red] [white]{href}[/white]")
    if verbose:
        console.print(f"[bold cyan]\nResumen: {v3_count} v3 válidos, {v2_count} v2 obsoletos, {otros} otros. {len(enlaces)} únicos.[/bold cyan]")
    return list(enlaces)

# === EXTRACCIÓN DE METADATOS Y ESTADO ===
def extraer_metadatos(onion_url, usar_tor=True):
    headers = {"User-Agent": random.choice(USER_AGENTS)}
    proxies = {}
    if usar_tor:
        host, port = TOR_SOCKS.split(":")
        proxies = {
            "http": f"socks5h://{host}:{port}",
            "https": f"socks5h://{host}:{port}"
        }
    try:
        r = requests.get(onion_url, headers=headers, timeout=15, proxies=proxies)
        soup = BeautifulSoup(r.text, "html.parser")
        title = soup.title.string.strip() if soup.title else "(sin título)"
        return {"url": onion_url, "title": title, "status": r.status_code}
    except Exception as e:
        return {"url": onion_url, "title": "(no accesible)", "status": str(e)}

# === DORKS/FUZZER ===
DORKS = [
    "admin", "login", "config", "backup", "robots.txt", ".git", "upload", "panel", "cpanel", "shell", "passwd", "secret", "hidden", "private", "db", "database", "dump", "leak", "test", "root", ".env",
    ".bak", ".old", ".zip", ".tar.gz", ".php", ".asp", ".jsp", ".cgi", "webadmin", "dashboard", "console", "server-status", "error.log", "access.log", ".htpasswd", ".htaccess", ".svn", ".DS_Store", ".idea", ".vscode", ".ssh", ".aws", "phpinfo.php", "info.php", "debug", "staging", "dev", "test", "tmp", "temp", "logs", "log", "dump.sql", "db.sql", "db_backup", "old", "archive", "backup.zip", "backup.tar.gz"
]

def fuzzear_onion(base_url, usar_tor=True, verbose=True):
    hallazgos = []
    for dork in DORKS:
        url = base_url.rstrip('/') + '/' + dork
        meta = extraer_metadatos(url, usar_tor=usar_tor)
        if meta['status'] and (str(meta['status']).startswith('2') or str(meta['status']).startswith('3')):
            console.print(f"[bold magenta][dork][+] {url}[/bold magenta] | [green]{meta['title']}[/green] | [yellow]status: {meta['status']}[/yellow]")
            hallazgos.append(meta)
        time.sleep(random.uniform(0.2, 0.5))
    return hallazgos

# === ESCANEO DE PUERTOS (tipo nmap) ===
def escanear_puertos(onion_url, puertos=None, timeout=3):
    if puertos is None:
        puertos = [80, 443, 8080, 8000, 8443, 8888, 22, 21, 3306, 53, 25, 110, 143, 465, 993, 995, 6667, 23, 137, 139, 445, 5900, 6660, 6669, 7000, 9001, 9050, 9051, 5000, 1234, 31337]
    host = onion_url.replace('http://', '').replace('https://', '').split('/')[0]
    abiertos = []
    for puerto in puertos:
        try:
            s = socks.socksocket()
            s.set_proxy(socks.SOCKS5, TOR_SOCKS.split(':')[0], int(TOR_SOCKS.split(':')[1]))
            s.settimeout(timeout)
            s.connect((host, puerto))
            abiertos.append(puerto)
            s.close()
            console.print(f"[bold blue][port][+] {host}:{puerto} abierto[/bold blue]")
        except Exception:
            pass
    return abiertos

# === BANNER GRABBER ===
def obtener_banner(onion_url, puerto, timeout=4):
    host = onion_url.replace('http://', '').replace('https://', '').split('/')[0]
    try:
        s = socks.socksocket()
        s.set_proxy(socks.SOCKS5, TOR_SOCKS.split(':')[0], int(TOR_SOCKS.split(':')[1]))
        s.settimeout(timeout)
        s.connect((host, puerto))
        s.sendall(b'HEAD / HTTP/1.0\r\nHost: ' + host.encode() + b'\r\n\r\n')
        data = s.recv(256)
        banner = data.decode(errors='ignore').split('\r\n')[0]
        s.close()
        console.print(f"[bold yellow][banner][{puerto}] {banner}[/bold yellow]")
        return banner
    except Exception:
        return None

# === DETECCIÓN DE TECNOLOGÍAS ===
def detectar_tecnologias(meta):
    techs = []
    if 'php' in meta['title'].lower(): techs.append('PHP')
    if 'nginx' in meta['title'].lower(): techs.append('nginx')
    if 'apache' in meta['title'].lower(): techs.append('Apache')
    if 'django' in meta['title'].lower(): techs.append('Django')
    if 'wordpress' in meta['title'].lower(): techs.append('WordPress')
    if 'ftp' in meta['title'].lower(): techs.append('FTP')
    if 'ssh' in meta['title'].lower(): techs.append('SSH')
    return techs

# === BUSCADORES ONION MODULARES ===
def buscar_ahmia(termino, usar_tor=True):
    url = f"https://ahmia.fi/search/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_darksearch(termino, usar_tor=True):
    url = f"https://darksearch.io/search/?query={termino}"
    return _buscar_en_motor(url, usar_tor)

# === MÁS MOTORES ONION (AJUSTADOS Y SEGUROS) ===
def buscar_phobos(termino, usar_tor=True):
    url = f"https://phobos.torsearch.com/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_onionland(termino, usar_tor=True):
    url = f"https://onionlandsearchengine.com/search/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_tordex(termino, usar_tor=True):
    url = f"http://tordex7xk5h3xq72.onion/search/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_haystak(termino, usar_tor=True):
    url = f"http://haystakvxad7wbk5.onion/search?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_torch(termino, usar_tor=True):
    url = f"http://cnkj6n6jraib3v7w.onion/search?query={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_notevil(termino, usar_tor=True):
    url = f"http://hss3uro2hsxfogfq.onion/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_kilos(termino, usar_tor=True):
    url = f"http://kilos7encrkxjfy.onion/search?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_recon(termino, usar_tor=True):
    url = f"http://reconponydonugup.onion/search/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_candle(termino, usar_tor=True):
    url = f"http://gjobqjj7wyczbqie.onion/search/{termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_deepsearch(termino, usar_tor=True):
    url = f"http://deepsearch6f6p7.onion/search/?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_onionsearch(termino, usar_tor=True):
    url = f"https://onionsearchengine.com/search.php?q={termino}"
    return _buscar_en_motor(url, usar_tor)

def buscar_darknetlive(termino, usar_tor=True):
    url = f"https://darknetlive.com/?s={termino}"
    return _buscar_en_motor(url, usar_tor)

# Función interna para parsear resultados de motores onion
def _buscar_en_motor(url, usar_tor):
    headers = {"User-Agent": random.choice(USER_AGENTS)}
    proxies = {}
    if usar_tor:
        host, port = TOR_SOCKS.split(":")
        proxies = {"http": f"socks5h://{host}:{port}", "https": f"socks5h://{host}:{port}"}
    try:
        r = requests.get(url, headers=headers, timeout=20, proxies=proxies)
        soup = BeautifulSoup(r.text, "html.parser")
        enlaces = set()
        for a in soup.find_all("a", href=True):
            href = a["href"].strip()
            if es_onion_v3(href):
                enlaces.add(href)
        return list(enlaces)
    except Exception:
        return []

# === FILTRADO Y DEDUPLICADO ===
def filtrar_onions(lista):
    # Solo .onion v3, sin duplicados ni falsos positivos
    return sorted(set([x for x in lista if es_onion_v3(x)]))

# === DETECCIÓN DE CAPTCHA/BLOQUEO ===
def es_captcha_o_bloqueo(html):
    patrones = ["captcha", "cloudflare", "are you human", "protection", "blocked", "verify"]
    return any(pat in html.lower() for pat in patrones)

# === MAIN ===
def main():
    try:
        keyword = console.input("[bold cyan]Palabra(s) clave .onion (separa por coma): [/bold cyan]").strip()
        if not keyword:
            console.print("[bold red][!] No ingresaste ninguna palabra clave.[/bold red]")
            sys.exit(1)
        keywords = [k.strip() for k in keyword.split(",") if k.strip()]
        usar_tor = True
        all_hallazgos = []
        motores = [buscar_ahmia, buscar_darksearch, buscar_phobos, buscar_onionland, buscar_tordex, buscar_haystak, buscar_torch, buscar_notevil, buscar_kilos, buscar_recon, buscar_candle, buscar_deepsearch, buscar_onionsearch, buscar_darknetlive]
        console.print(f"\n[bold yellow][~] Buscando en {len(motores)} motores onion... (Ctrl+C para salir y guardar)[/bold yellow]\n")
        encontrados = set()
        for motor in motores:
            for kw in keywords:
                try:
                    resultados = motor(kw, usar_tor=usar_tor)
                except Exception as e:
                    console.print(f"[bold red][!] Motor {motor.__name__} falló: {e}[/bold red]")
                    continue
                for r in filtrar_onions(resultados):
                    if r in encontrados:
                        continue
                    encontrados.add(r)
                    meta = extraer_metadatos(r, usar_tor=usar_tor)
                    if meta['title'] and es_captcha_o_bloqueo(meta['title']):
                        console.print(f"[bold red][captcha][!] {r} parece estar protegido o bloqueado[/bold red]")
                        continue
                    console.print(f"[bold green][buscador][+] {meta['url']}[/bold green] | [white]{meta['title']}[/white] | [yellow]status: {meta['status']}[/yellow]")
                    all_hallazgos.append(meta)
                    puertos_abiertos = escanear_puertos(r)
                    for p in puertos_abiertos:
                        banner = obtener_banner(r, p)
                    dork_hallazgos = fuzzear_onion(r, usar_tor=usar_tor, verbose=False)
                    for dork in dork_hallazgos:
                        console.print(f"[bold magenta][dork][+] {dork['url']}[/bold magenta] | [green]{dork['title']}[/green] | [yellow]status: {dork['status']}[/yellow]")
                    all_hallazgos.extend(dork_hallazgos)
                    techs = detectar_tecnologias(meta)
                    if techs:
                        console.print(f"[bold cyan][tech][{meta['url']}] Tecnologías detectadas: {', '.join(techs)}[/bold cyan]")
                    time.sleep(random.uniform(0.5, 1.2))
        # Resumen y guardado
        total = len(all_hallazgos)
        paneles = len([h for h in all_hallazgos if any(x in h['url'] for x in ['admin','panel','cpanel','dashboard','console'])])
        abiertos = len([h for h in all_hallazgos if str(h['status']).startswith('2') or str(h['status']).startswith('3')])
        console.print(f"\n[bold cyan]Resumen: {total} hallazgos, {paneles} posibles paneles, {abiertos} accesibles.[/bold cyan]")
        table = Table(title="Resultados .onion", show_lines=True, style="bold green")
        table.add_column("URL", style="white", overflow="fold")
        table.add_column("Título", style="green")
        table.add_column("Status", style="yellow")
        for meta in all_hallazgos:
            table.add_row(meta['url'], meta['title'], str(meta['status']))
        console.print(table)
        with open("resultados_onion_fuzz.txt", "w") as f:
            for meta in all_hallazgos:
                f.write(f"{meta['url']} | {meta['title']} | status: {meta['status']}\n")
        console.print("\n[bold magenta]📁 Resultados guardados en 'resultados_onion_fuzz.txt'.[/bold magenta]")
    except KeyboardInterrupt:
        console.print("\n[bold red]⏹️ Proceso interrumpido por el usuario. Guardando resultados...[/bold red]")
        with open("resultados_onion_fuzz.txt", "w") as f:
            for meta in all_hallazgos:
                f.write(f"{meta['url']} | {meta['title']} | status: {meta['status']}\n")
        console.print("\n[bold magenta]📁 Resultados guardados en 'resultados_onion_fuzz.txt'.[/bold magenta]")
        sys.exit(0)
