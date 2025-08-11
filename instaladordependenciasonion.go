package main

import (
	"fmt"
	"os/exec"
	"strings"
)

func main() {
	fmt.Println("# Script para instalar dependencias de buscaonion (Go version)")
	fmt.Println("Dependencias Python:")
	fmt.Println("  sudo apt-get update && sudo apt-get install -y python3 python3-pip")
	fmt.Println("  pip3 install requests[socks] beautifulsoup4 rich PySocks")
	fmt.Println("Dependencias Go:")
	fmt.Println("  go get github.com/PuerkitoBio/goquery")
	fmt.Print("¿Deseas instalar dependencias Go automáticamente? (s/n): ")
	var resp string
	fmt.Scanln(&resp)
	if strings.ToLower(resp) == "s" {
		cmd := exec.Command("go", "mod", "init", "martillazonion-go")
		cmd.Run()
		cmd = exec.Command("go", "mod", "tidy")
		cmd.Run()
		fmt.Println("[+] Dependencias Go instaladas automáticamente.")
	}
	fmt.Println("[+] Dependencias listadas. Ya puedes usar buscaonion en Go.")
}
