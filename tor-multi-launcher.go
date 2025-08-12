package main

import (
	"flag"
	"fmt"
	"os"
)

func cargarRangoDesdeArchivo(path string) (int, int) {
	f, err := os.Open(path)
	if err != nil {
		return 9055, 9060
	}
	defer f.Close()
	var inicio, fin int
	fmt.Fscanf(f, "%d %d", &inicio, &fin)
	if inicio == 0 || fin == 0 {
		return 9055, 9060
	}
	return inicio, fin
}

func main() {
	inicio := flag.Int("inicio", 9055, "Puerto inicial")
	fin := flag.Int("fin", 9060, "Puerto final")
	var configFile string
	flag.StringVar(&configFile, "config", "", "Archivo con rango de puertos")
	flag.Parse()
	if configFile != "" {
		*inicio, *fin = cargarRangoDesdeArchivo(configFile)
	}
	fmt.Printf("[+] Iniciando instancias Tor en puertos %d a %d...\n", *inicio, *fin)
	total := 0
	for port := *inicio; port <= *fin; port++ {
		dir := fmt.Sprintf("/var/lib/tor%d", port)
		log := fmt.Sprintf("/var/log/tor%d.log", port)
		conf := fmt.Sprintf("/etc/tor/torrc%d", port)
		fmt.Printf(" - Crear directorio: sudo mkdir -p %s\n", dir)
		fmt.Printf(" - Cambiar permisos: sudo chown debian-tor:debian-tor %s && sudo chmod 700 %s\n", dir, dir)
		fmt.Printf(" - Crear config: sudo bash -c 'cat > %s <<EOF\\nSocksPort 127.0.0.1:%d\\nDataDirectory %s\\nLog notice file %s\\nEOF'\n", conf, port, dir, log)
		fmt.Printf(" - Lanzar Tor: sudo -u debian-tor tor -f %s &\n", conf)
		total++
	}
	fmt.Printf("[✓] %d instancias de Tor han sido lanzadas (simulado).\n", total)
}
