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
	fmt.Printf("[!] Deteniendo instancias Tor en puertos %d a %d...\n", *inicio, *fin)
	total := 0
	for port := *inicio; port <= *fin; port++ {
		conf := fmt.Sprintf("/etc/tor/torrc%d", port)
		log := fmt.Sprintf("/var/log/tor%d.log", port)
		dir := fmt.Sprintf("/var/lib/tor%d", port)
		fmt.Printf(" - Terminando procesos en puerto %d: sudo pkill -f \"tor.*%s\"\n", port, conf)
		// Opcional: eliminar archivos
		// fmt.Printf("   - Borrando configuración y datos: sudo rm -f \"%s\" \"%s\" && sudo rm -rf \"%s\"\n", conf, log, dir)
		total++
	}
	fmt.Printf("[✓] %d instancias Tor han sido detenidas (simulado).\n", total)
}
