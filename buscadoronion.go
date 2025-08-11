package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func mostrarBanner() {
	banner := `
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
`
	fmt.Println(banner)
}

func buscarOnion(termino string) ([]string, error) {
	fmt.Printf("\n🔍 Buscando sitios .onion relacionados con: %s\n", termino)
	url := fmt.Sprintf("https://ahmia.fi/search/?q=%s", termino)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("[!] Error al conectar con Ahmia: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("[!] Código de estado HTTP inesperado: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	enlaces := make(map[string]struct{})
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists && strings.Contains(href, ".onion") {
			enlaces[strings.TrimSpace(href)] = struct{}{}
		}
	})

	var resultados []string
	for enlace := range enlaces {
		resultados = append(resultados, enlace)
	}
	return resultados, nil
}

// BuscarEnArchivo busca una palabra clave en un archivo de texto y muestra las líneas que la contienen.
func BuscarEnArchivo(rutaArchivo string, palabraClave string) error {
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return err
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	lineaNum := 1
	for scanner.Scan() {
		linea := scanner.Text()
		if strings.Contains(linea, palabraClave) {
			fmt.Printf("Línea %d: %s\n", lineaNum, linea)
		}
		lineaNum++
	}
	return scanner.Err()
}

// Validadores de dominios .onion v2/v3
var (
	onionV3Regex = regexp.MustCompile(`[a-z2-7]{56}\.onion`)
	onionV2Regex = regexp.MustCompile(`[a-z2-7]{16}\.onion`)
)

// EsOnionV3 verifica si el string es un dominio .onion v3
func EsOnionV3(url string) bool {
	return onionV3Regex.MatchString(url)
}

// EsOnionV2 verifica si el string es un dominio .onion v2
func EsOnionV2(url string) bool {
	return onionV2Regex.MatchString(url)
}

// ExtraerTitulo obtiene el título de una página .onion
func ExtraerTitulo(url string) string {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "(error)"
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return "(no accesible)"
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "(sin título)"
	}
	title := strings.TrimSpace(doc.Find("title").Text())
	if title == "" {
		return "(sin título)"
	}
	return title
}

// DetectarTecnologias busca tecnologías comunes en el título
func DetectarTecnologias(title string) []string {
	techs := []string{}
	t := strings.ToLower(title)
	if strings.Contains(t, "php") {
		techs = append(techs, "PHP")
	}
	if strings.Contains(t, "nginx") {
		techs = append(techs, "nginx")
	}
	if strings.Contains(t, "apache") {
		techs = append(techs, "Apache")
	}
	if strings.Contains(t, "django") {
		techs = append(techs, "Django")
	}
	if strings.Contains(t, "wordpress") {
		techs = append(techs, "WordPress")
	}
	if strings.Contains(t, "ftp") {
		techs = append(techs, "FTP")
	}
	if strings.Contains(t, "ssh") {
		techs = append(techs, "SSH")
	}
	return techs
}

// FiltrarYDeduplicarOnions filtra solo .onion v3 y elimina duplicados
func FiltrarYDeduplicarOnions(lista []string) []string {
	uniq := make(map[string]struct{})
	for _, x := range lista {
		if EsOnionV3(x) {
			uniq[x] = struct{}{}
		}
	}
	var res []string
	for k := range uniq {
		res = append(res, k)
	}
	return res
}

// buscarEnMotores permite buscar en varios motores (solo Ahmia implementado)
func buscarEnMotores(termino string) []string {
	// Aquí puedes agregar más motores en el futuro
	resultados, _ := buscarOnion(termino)
	return resultados
}

// mostrarResumenTecnologias muestra tecnologías detectadas en los resultados
func mostrarResumenTecnologias(resultados map[string][]string) {
	fmt.Println("\n🔬 Tecnologías detectadas por palabra clave:")
	for k, enlaces := range resultados {
		tecs := make(map[string]int)
		for _, url := range enlaces {
			titulo := ExtraerTitulo(url)
			for _, tech := range DetectarTecnologias(titulo) {
				tecs[tech]++
			}
		}
		if len(tecs) == 0 {
			fmt.Printf("  - '%s': (ninguna detectada)\n", k)
		} else {
			fmt.Printf("  - '%s': ", k)
			for tech, count := range tecs {
				fmt.Printf("%s(%d) ", tech, count)
			}
			fmt.Println()
		}
	}
}

// BuscarEnArchivo busca una palabra clave en un archivo de texto y muestra las líneas que la contienen.
func BuscarEnArchivo(rutaArchivo string, palabraClave string) error {
	archivo, err := os.Open(rutaArchivo)
	if err != nil {
		return err
	}
	defer archivo.Close()

	scanner := bufio.NewScanner(archivo)
	lineaNum := 1
	for scanner.Scan() {
		linea := scanner.Text()
		if strings.Contains(linea, palabraClave) {
			fmt.Printf("Línea %d: %s\n", lineaNum, linea)
		}
		lineaNum++
	}
	return scanner.Err()
}

// buscarMultiplesKeywords permite buscar varios términos y agrupa los resultados.
func buscarMultiplesKeywords(terminos []string) map[string][]string {
	resultados := make(map[string][]string)
	for _, termino := range terminos {
		enlaces := buscarEnMotores(termino)
		enlaces = FiltrarYDeduplicarOnions(enlaces)
		resultados[termino] = enlaces
	}
	return resultados
}

// leerKeywordsDesdeArchivo lee palabras clave desde un archivo de texto (una por línea).
func leerKeywordsDesdeArchivo(ruta string) ([]string, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var keywords []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		kw := strings.TrimSpace(scanner.Text())
		if kw != "" {
			keywords = append(keywords, kw)
		}
	}
	return keywords, scanner.Err()
}

// mostrarEstadisticas muestra un resumen de los resultados encontrados.
func mostrarEstadisticas(resultados map[string][]string) {
	fmt.Println("\n📊 Estadísticas de búsqueda:")
	total := 0
	for k, v := range resultados {
		fmt.Printf("  - '%s': %d enlaces\n", k, len(v))
		total += len(v)
	}
	fmt.Printf("  Total de enlaces encontrados: %d\n", total)
}

// ExportarResultadosJSON exporta los resultados a un archivo JSON.
func ExportarResultadosJSON(nombreArchivo string, resultados map[string][]string) error {
	data := make(map[string][]map[string]interface{})
	for k, enlaces := range resultados {
		var lista []map[string]interface{}
		for _, url := range enlaces {
			titulo := ExtraerTitulo(url)
			tecs := DetectarTecnologias(titulo)
			lista = append(lista, map[string]interface{}{
				"url":    url,
				"titulo": titulo,
				"tecs":   tecs,
			})
		}
		data[k] = lista
	}
	f, err := os.Create(nombreArchivo)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// MostrarSoloV3 imprime solo los dominios .onion v3 de los resultados.
func MostrarSoloV3(resultados map[string][]string) {
	fmt.Println("\n🔎 Solo dominios .onion v3 encontrados:")
	for k, enlaces := range resultados {
		fmt.Printf("  [%s]:\n", k)
		for _, url := range enlaces {
			if EsOnionV3(url) {
				fmt.Println("   -", url)
			}
		}
	}
}

// MostrarSoloV2 imprime solo los dominios .onion v2 de los resultados.
func MostrarSoloV2(resultados map[string][]string) {
	fmt.Println("\n🔎 Solo dominios .onion v2 encontrados:")
	for k, enlaces := range resultados {
		fmt.Printf("  [%s]:\n", k)
		for _, url := range enlaces {
			if EsOnionV2(url) {
				fmt.Println("   -", url)
			}
		}
	}
}

// MostrarEstadisticasTitulos muestra estadísticas de longitud de títulos.
func MostrarEstadisticasTitulos(resultados map[string][]string) {
	fmt.Println("\n📏 Estadísticas de longitud de títulos:")
	for k, enlaces := range resultados {
		var total, count int
		for _, url := range enlaces {
			titulo := ExtraerTitulo(url)
			total += len(titulo)
			count++
		}
		if count > 0 {
			fmt.Printf("  [%s]: Promedio de longitud de título: %.2f\n", k, float64(total)/float64(count))
		}
	}
}

// BuscarTecnologiaEnResultados busca una tecnología específica en los títulos.
func BuscarTecnologiaEnResultados(resultados map[string][]string, tecnologia string) {
	fmt.Printf("\n🔬 Buscando tecnología '%s' en títulos:\n", tecnologia)
	tec := strings.ToLower(tecnologia)
	for k, enlaces := range resultados {
		for _, url := range enlaces {
			titulo := ExtraerTitulo(url)
			if strings.Contains(strings.ToLower(titulo), tec) {
				fmt.Printf("  [%s] %s | %s\n", k, url, titulo)
			}
		}
	}
}

// Resultados globales y mutex para acceso concurrente
var (
	ultimosResultados     = make(map[string][]string)
	ultimosResultadosLock sync.RWMutex
)

// SetUltimosResultados actualiza los resultados globales
func SetUltimosResultados(res map[string][]string) {
	ultimosResultadosLock.Lock()
	defer ultimosResultadosLock.Unlock()
	ultimosResultados = res
}

// GetUltimosResultados obtiene una copia de los resultados globales
func GetUltimosResultados() map[string][]string {
	ultimosResultadosLock.RLock()
	defer ultimosResultadosLock.RUnlock()
	copia := make(map[string][]string)
	for k, v := range ultimosResultados {
		copia[k] = append([]string{}, v...)
	}
	return copia
}

// LimpiarUltimosResultados borra los resultados globales
func LimpiarUltimosResultados() {
	ultimosResultadosLock.Lock()
	defer ultimosResultadosLock.Unlock()
	ultimosResultados = make(map[string][]string)
}

// HayResultadosEnMemoria indica si hay resultados guardados
func HayResultadosEnMemoria() bool {
	ultimosResultadosLock.RLock()
	defer ultimosResultadosLock.RUnlock()
	return len(ultimosResultados) > 0
}

// menuInteractivo permite elegir el modo de búsqueda.
func menuInteractivo() {
	mostrarBanner()
	fmt.Println("🧅 Buscador interactivo de sitios .onion\n")
	fmt.Println("1. Buscar una palabra clave")
	fmt.Println("2. Buscar varias palabras clave (separadas por coma)")
	fmt.Println("3. Buscar palabras clave desde archivo")
	fmt.Println("4. Buscar en archivo local de resultados")
	fmt.Println("5. Limpiar resultados en memoria")
	fmt.Println("0. Salir")
	fmt.Print("\nSelecciona una opción: ")

	reader := bufio.NewReader(os.Stdin)
	opcion, _ := reader.ReadString('\n')
	opcion = strings.TrimSpace(opcion)

	switch opcion {
	case "1":
		fmt.Print("Ingresa una palabra clave: ")
		keyword, _ := reader.ReadString('\n')
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			fmt.Println("[!] No ingresaste ninguna palabra clave.")
			return
		}
		resultados := buscarEnMotores(keyword)
		resultados = FiltrarYDeduplicarOnions(resultados)
		if len(resultados) == 0 {
			fmt.Println("⚠️ No se encontraron resultados.")
		} else {
			fmt.Printf("\n✅ Se encontraron %d enlaces .onion:\n\n", len(resultados))
			for _, r := range resultados {
				fmt.Println("   -", r)
				titulo := ExtraerTitulo(r)
				fmt.Printf("     Título: %s\n", titulo)
				tecs := DetectarTecnologias(titulo)
				if len(tecs) > 0 {
					fmt.Printf("     Tecnologías: %s\n", strings.Join(tecs, ", "))
				}
			}
			f, err := os.Create("resultados_onion.txt")
			if err != nil {
				fmt.Println("[!] No se pudo guardar el archivo de resultados.")
				return
			}
			defer f.Close()
			for _, r := range resultados {
				f.WriteString(r + "\n")
			}
			fmt.Println("\n📁 Resultados guardados en 'resultados_onion.txt'.")
			// Actualiza resultados globales
			SetUltimosResultados(map[string][]string{"busqueda": resultados})
		}
	case "2":
		fmt.Print("Ingresa palabras clave separadas por coma: ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		terminos := strings.Split(line, ",")
		for i := range terminos {
			terminos[i] = strings.TrimSpace(terminos[i])
		}
		resultados := buscarMultiplesKeywords(terminos)
		mostrarEstadisticas(resultados)
		mostrarResumenTecnologias(resultados)
		f, err := os.Create("resultados_onion_multi.txt")
		if err == nil {
			defer f.Close()
			for k, enlaces := range resultados {
				for _, r := range enlaces {
					f.WriteString(fmt.Sprintf("[%s] %s\n", k, r))
				}
			}
			fmt.Println("\n📁 Resultados guardados en 'resultados_onion_multi.txt'.")
		}
		SetUltimosResultados(resultados)
	case "3":
		fmt.Print("Ruta del archivo de palabras clave: ")
		ruta, _ := reader.ReadString('\n')
		ruta = strings.TrimSpace(ruta)
		terminos, err := leerKeywordsDesdeArchivo(ruta)
		if err != nil {
			fmt.Println("[!] No se pudo leer el archivo:", err)
			return
		}
		resultados := buscarMultiplesKeywords(terminos)
		mostrarEstadisticas(resultados)
		mostrarResumenTecnologias(resultados)
		f, err := os.Create("resultados_onion_archivo.txt")
		if err == nil {
			defer f.Close()
			for k, enlaces := range resultados {
				for _, r := range enlaces {
					f.WriteString(fmt.Sprintf("[%s] %s\n", k, r))
				}
			}
			fmt.Println("\n📁 Resultados guardados en 'resultados_onion_archivo.txt'.")
		}
		SetUltimosResultados(resultados)
	case "4":
		fmt.Print("Ruta del archivo de resultados: ")
		ruta, _ := reader.ReadString('\n')
		ruta = strings.TrimSpace(ruta)
		fmt.Print("Palabra clave a buscar: ")
		clave, _ := reader.ReadString('\n')
		clave = strings.TrimSpace(clave)
		err := BuscarEnArchivo(ruta, clave)
		if err != nil {
			fmt.Println("[!] Error al buscar en archivo:", err)
		}
	case "5":
		LimpiarUltimosResultados()
		fmt.Println("Resultados en memoria limpiados.")
	case "0":
		fmt.Println("¡Hasta luego!")
		os.Exit(0)
	default:
		fmt.Println("[!] Opción no válida.")
	}
}

func main() {
	for {
		menuInteractivo()
		fmt.Println("\n----------------------------------------\n")
	}
}
