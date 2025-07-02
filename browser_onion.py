#!/usr/bin/env python3
import subprocess
import sys

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

import time

def mostrar_banner():
    banner = """
 ▄▄▄▄    █    ██   ██████  ▄████▄   ▄▄▄      ▓█████▄  ▒█████   ██▀███        ▒█████   ███▄    █  ██▓ ▒█████   ███▄    █ 
▓█████▄  ██  ▓██▒▒██    ▒ ▒██▀ ▀█  ▒████▄    ▒██▀ ██▌▒██▒  ██▒▓██ ▒ ██▒     ▒██▒  ██▒ ██ ▀█   █ ▓██▒▒██▒  ██▒ ██ ▀█   █ 
▒██▒ ▄██▓██  ▒██░░ ▓██▄   ▒▓█    ▄ ▒██  ▀█▄  ░██   █▌▒██░  ██▒▓██ ░▄█ ▒     ▒██░  ██▒▓██  ▀█ ██▒▒██▒▒██░  ██▒▓██  ▀█ ██▒
▒██░█▀  ▓▓█  ░██░  ▒   ██▒▒▓▓▄ ▄██▒░██▄▄▄▄██ ░▓█▄   ▌▒██   ██░▒██▀▀█▄       ▒██   ██░▓██▒  ▐▌██▒░██░▒██   ██░▓██▒  ▐▌██▒
░▓█  ▀█▓▒▒█████▓ ▒██████▒▒▒ ▓███▀ ░ ▓█   ▓██▒░▒████▓ ░ ████▓▒░░██▓ ▒██▒ ██▓ ░ ████▓▒░▒██░   ▓██░░██░░ ████▓▒░▒██░   ▓██░
░▒▓███▀▒░▒▓▒ ▒ ▒ ▒ ▒▓▒ ▒ ░░ ░▒ ▒  ░ ▒▒   ▓▒█░ ▒▒▓  ▒ ░ ▒░▒░▒░ ░ ▒▓ ░▒▓░ ▒▓▒ ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ ░▓  ░ ▒░▒░▒░ ░ ▒░   ▒ ▒ 
▒░▒   ░ ░░▒░ ░ ░ ░ ░▒  ░ ░  ░  ▒     ▒   ▒▒ ░ ░ ▒  ▒   ░ ▒ ▒░   ░▒ ░ ▒░ ░▒    ░ ▒ ▒░ ░ ░░   ░ ▒░ ▒ ░  ░ ▒ ▒░ ░ ░░   ░ ▒░
 ░    ░  ░░░ ░ ░ ░  ░  ░  ░          ░   ▒    ░ ░  ░ ░ ░ ░ ▒    ░░   ░  ░   ░ ░ ░ ▒     ░   ░ ░  ▒ ░░ ░ ░ ▒     ░   ░ ░ 
 ░         ░           ░  ░ ░            ░  ░   ░        ░ ░     ░       ░      ░ ░           ░  ░      ░ ░           ░ 
      ░                   ░                   ░                          ░                                              
"""
    print(banner)

def buscar_onion(termino):
    print(f"\n🔍 Buscando sitios .onion relacionados con: {termino}")
    url = f"https://ahmia.fi/search/?q={termino}"
    headers = {
        "User-Agent": "Mozilla/5.0"
    }

    try:
        respuesta = requests.get(url, headers=headers, timeout=15)
        respuesta.raise_for_status()
    except Exception as e:
        print(f"[!] Error al conectar con Ahmia: {e}")
        return []

    soup = BeautifulSoup(respuesta.text, "html.parser")
    enlaces = []

    for a in soup.find_all("a", href=True):
        href = a["href"]
        if ".onion" in href:
            enlaces.append(href.strip())

    return list(set(enlaces))

def main():
    mostrar_banner()
    print("🧅 Buscador interactivo de sitios .onion\n")
    keyword = input("Ingresa una palabra clave (ej. market, forum, bitcoin): ").strip()

    if not keyword:
        print("[!] No ingresaste ninguna palabra clave.")
        sys.exit(1)

    resultados = buscar_onion(keyword)

    if not resultados:
        print("⚠️ No se encontraron resultados.")
    else:
        print(f"\n✅ Se encontraron {len(resultados)} enlaces .onion:\n")
        for r in resultados:
            print("   -", r)

        with open("resultados_onion.txt", "w") as f:
            for r in resultados:
                f.write(r + "\n")
        print("\n📁 Resultados guardados en 'resultados_onion.txt'.")

if __name__ == "__main__":
    main()
