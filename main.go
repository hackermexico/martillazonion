package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	for {
		fmt.Println("Bienvenido a martillazonion-go")
		fmt.Println("1. Menú buscador interactivo")
		fmt.Println("2. Exportar últimos resultados a JSON")
		fmt.Println("3. Mostrar solo dominios .onion v3")
		fmt.Println("4. Mostrar solo dominios .onion v2")
		fmt.Println("5. Estadísticas de longitud de títulos")
		fmt.Println("6. Buscar tecnología en títulos")
		fmt.Println("0. Salir")
		fmt.Print("Selecciona una opción: ")

		reader := bufio.NewReader(os.Stdin)
		op, _ := reader.ReadString('\n')
		op = strings.TrimSpace(op)

		switch op {
		case "1":
			menuInteractivo()
		case "2":
			if !HayResultadosEnMemoria() {
				fmt.Println("No hay resultados previos en memoria. Usa el menú interactivo primero.")
				break
			}
			fmt.Print("Nombre de archivo JSON a exportar: ")
			nombre, _ := reader.ReadString('\n')
			nombre = strings.TrimSpace(nombre)
			err := ExportarResultadosJSON(nombre, GetUltimosResultados())
			if err != nil {
				fmt.Println("Error exportando a JSON:", err)
			} else {
				fmt.Println("Exportado correctamente a", nombre)
			}
		case "3":
			if !HayResultadosEnMemoria() {
				fmt.Println("No hay resultados previos en memoria. Usa el menú interactivo primero.")
				break
			}
			MostrarSoloV3(GetUltimosResultados())
		case "4":
			if !HayResultadosEnMemoria() {
				fmt.Println("No hay resultados previos en memoria. Usa el menú interactivo primero.")
				break
			}
			MostrarSoloV2(GetUltimosResultados())
		case "5":
			if !HayResultadosEnMemoria() {
				fmt.Println("No hay resultados previos en memoria. Usa el menú interactivo primero.")
				break
			}
			MostrarEstadisticasTitulos(GetUltimosResultados())
		case "6":
			if !HayResultadosEnMemoria() {
				fmt.Println("No hay resultados previos en memoria. Usa el menú interactivo primero.")
				break
			}
			fmt.Print("Tecnología a buscar: ")
			tec, _ := reader.ReadString('\n')
			tec = strings.TrimSpace(tec)
			BuscarTecnologiaEnResultados(GetUltimosResultados(), tec)
		case "0":
			fmt.Println("¡Hasta luego!")
			os.Exit(0)
		default:
			fmt.Println("Opción no válida.")
		}
		fmt.Println()
	}
}
