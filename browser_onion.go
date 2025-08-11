package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
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

func buscarOnion(termino string) []string {
	fmt.Printf("\n🔍 Buscando sitios .onion relacionados con: %s\n", termino)
	url := fmt.Sprintf("https://ahmia.fi/search/?q=%s", termino)
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[!] Error al conectar con Ahmia: %v\n", err)
		return nil
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("[!] Error al parsear HTML: %v\n", err)
		return nil
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
	return resultados
}

func buscarDesdeArchivo(path string) map[string]struct{} {
	resultadosUnicos := make(map[string]struct{})
	f, err := os.Open(path)
	if err != nil {
		fmt.Println("No se pudo abrir el archivo:", err)
		return resultadosUnicos
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		kw := strings.TrimSpace(scanner.Text())
		if kw == "" {
			continue
		}
		resultados := buscarOnion(kw)
		for _, r := range resultados {
			resultadosUnicos[r] = struct{}{}
		}
	}
	return resultadosUnicos
}

func exportarOnionLinksJSON(links map[string]struct{}, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	arr := make([]string, 0, len(links))
	for l := range links {
		arr = append(arr, l)
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(arr)
}

func main() {
	mostrarBanner()
	fmt.Println("🧅 Buscador interactivo de sitios .onion\n")
	fmt.Printf("[%s] ", time.Now().Format("15:04:05"))
	fmt.Print("¿Buscar desde archivo de palabras clave? (s/n): ")
	var resp string
	fmt.Scanln(&resp)
	var resultadosUnicos map[string]struct{}
	if strings.ToLower(resp) == "s" {
		fmt.Print("Ruta del archivo: ")
		var path string
		fmt.Scanln(&path)
		resultadosUnicos = buscarDesdeArchivo(path)
	} else {
		fmt.Printf("[%s] ", time.Now().Format("15:04:05"))
		fmt.Print("Ingresa palabra(s) clave (separa por coma): ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			fmt.Println("[!] No ingresaste ninguna palabra clave.")
			os.Exit(1)
		}
		keywords := strings.Split(line, ",")
		resultadosUnicos = make(map[string]struct{})
		for _, kw := range keywords {
			kw = strings.TrimSpace(kw)
			if kw == "" {
				continue
			}
			resultados := buscarOnion(kw)
			for _, r := range resultados {
				if strings.HasSuffix(r, ".onion") {
					resultadosUnicos[r] = struct{}{}
				}
			}
		}
	}
	if len(resultadosUnicos) == 0 {
		fmt.Println("⚠️ No se encontraron resultados.")
	} else {
		fmt.Printf("\n✅ Se encontraron %d enlaces .onion únicos:\n\n", len(resultadosUnicos))
		for r := range resultadosUnicos {
			fmt.Println("   -", r)
		}
		timestamp := time.Now().Format("20060102_150405")
		txtFile := fmt.Sprintf("resultados_onion_%s.txt", timestamp)
		jsonFile := fmt.Sprintf("resultados_onion_%s.json", timestamp)
		f, err := os.Create(txtFile)
		if err != nil {
			fmt.Println("[!] No se pudo guardar el archivo de resultados.")
			return
		}
		defer f.Close()
		for r := range resultadosUnicos {
			f.WriteString(r + "\n")
		}
		exportarOnionLinksJSON(resultadosUnicos, jsonFile)
		fmt.Printf("\n📁 Resultados guardados en '%s' y '%s'.\n", txtFile, jsonFile)
	}
}
