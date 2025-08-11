package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func mostrarBanner() {
	banner := `
 ░▒▓██████▓▒░░▒▓███████▓▒░ ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓████████▓▒░▒▓███████▓▒░        ░▒▓██████▓▒░░▒▓███████▓▒░░▒▓█▓▒░░▒▓██████▓▒░░▒▓███████▓▒░  
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░      ░▒▓███████▓▒░░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓██████▓▒░ ░▒▓███████▓▒░       ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░      ░▒▓█▓▒░      ░▒▓█▓▒░░▒▓█▓▒░▒▓██▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░ 
 ░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓█▓▒░░▒▓█████████████▓▒░░▒▓████████▓▒░▒▓████████▓▒░▒▓█▓▒░░▒▓█▓▒░▒▓██▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░▒▓█▓▒░░▒▓██████▓▒░░▒▓█▓▒░░▒▓█▓▒░ 
                                                                                                                                                                     
`
	fmt.Println(banner)
	fmt.Println("[!] Bienvenido al rastreador de sitios Onion\n")
}

func esUrlOnionValida(url string) bool {
	patron := regexp.MustCompile(`^http://[a-z2-7]{16,56}\.onion/?$`)
	return patron.MatchString(url)
}

func filtrarOnionLinks(links []string) []string {
	var onions []string
	re := regexp.MustCompile(`\.onion`)
	for _, l := range links {
		if re.MatchString(l) {
			onions = append(onions, l)
		}
	}
	return onions
}

func exportarOnionLinks(links []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	for _, l := range links {
		f.WriteString(l + "\n")
	}
	return nil
}

func exportarOnionLinksJSON(links []string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(links)
}

func rastrearOnion(url string) {
	fmt.Printf("\n[%s] Conectando a %s ...\n", time.Now().Format("15:04:05"), url)
	client := &http.Client{
		Timeout: 20 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("[-] Error al crear la petición:", err)
		os.Exit(1)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("[-] Error al conectar:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("[-] Error al parsear HTML:", err)
		os.Exit(1)
	}
	title := strings.TrimSpace(doc.Find("title").Text())
	if title == "" {
		title = "Sin título"
	}
	var links []string
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})
	onionLinks := filtrarOnionLinks(links)
	timestamp := time.Now().Format("20060102_150405")
	datos := map[string]interface{}{
		"url":     url,
		"titulo":  title,
		"enlaces": links,
		"fecha":   timestamp,
	}
	jsonFile := fmt.Sprintf("resultado_%s.json", timestamp)
	csvFile := fmt.Sprintf("resultado_%s.csv", timestamp)

	jfile, _ := os.Create(jsonFile)
	defer jfile.Close()
	json.NewEncoder(jfile).Encode(datos)

	cfile, _ := os.Create(csvFile)
	defer cfile.Close()
	writer := csv.NewWriter(cfile)
	defer writer.Flush()
	writer.Write([]string{"Enlaces encontrados"})
	for _, link := range links {
		writer.Write([]string{link})
	}

	fmt.Printf("\n[✔] Rastreo completado con éxito.\n[💾] Datos guardados en:\n - %s\n - %s\n", jsonFile, csvFile)
	fmt.Printf("[*] Enlaces .onion encontrados: %d\n", len(onionLinks))
	if len(onionLinks) > 0 {
		txtFile := fmt.Sprintf("enlaces_onion_%s.txt", timestamp)
		jsonFile := fmt.Sprintf("enlaces_onion_%s.json", timestamp)
		if err := exportarOnionLinks(onionLinks, txtFile); err == nil {
			fmt.Printf("[+] Enlaces .onion exportados a: %s\n", txtFile)
		}
		if err := exportarOnionLinksJSON(onionLinks, jsonFile); err == nil {
			fmt.Printf("[+] Enlaces .onion exportados a: %s\n", jsonFile)
		}
	}
}

func cargarObjetivosDesdeArchivo(archivo string) []string {
	var urls []string
	f, err := os.Open(archivo)
	if err != nil {
		fmt.Println("[-] Error al abrir el archivo:", err)
		return urls
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())
		if esUrlOnionValida(url) {
			urls = append(urls, url)
		} else {
			fmt.Printf("[-] URL no válida en el archivo: %s\n", url)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("[-] Error al leer el archivo:", err)
	}

	return urls
}

func rastrearMultiplesOnion(archivo string) {
	urls := cargarObjetivosDesdeArchivo(archivo)
	todosOnionLinks := make(map[string]struct{})
	for _, url := range urls {
		fmt.Printf("\n[+] Rastreo de: %s\n", url)
		rastrearOnion(url)
		// Extrae y acumula enlaces .onion
		client := &http.Client{Timeout: 20 * time.Second}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if exists && strings.Contains(href, ".onion") {
				todosOnionLinks[href] = struct{}{}
			}
		})
	}
	if len(todosOnionLinks) > 0 {
		f, _ := os.Create("enlaces_onion_global.txt")
		defer f.Close()
		for l := range todosOnionLinks {
			f.WriteString(l + "\n")
		}
		fmt.Printf("\n[+] Enlaces .onion globales exportados a: enlaces_onion_global.txt\n")
	}
	fmt.Printf("[*] Total de enlaces .onion únicos encontrados: %d\n", len(todosOnionLinks))
}

func main() {
	mostrarBanner()
	var url string
	fmt.Print("¿Deseas rastrear un archivo de URLs? (s/n): ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) == "s" {
		fmt.Print("Ruta del archivo: ")
		fmt.Scanln(&url)
		rastrearMultiplesOnion(url)
		return
	}
	for {
		fmt.Print("Introduce la URL .onion (incluyendo http://): ")
		fmt.Scanln(&url)
		if esUrlOnionValida(url) {
			break
		}
		fmt.Println("[-] URL .onion no válida. Intenta nuevamente.\n")
	}
	rastrearOnion(url)
}
