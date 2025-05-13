package handlers

import (
	"fmt"
	"time"
)

func MostrarRutaOptimizada() {
	fmt.Print("\033[H\033[2J")

	verde := "\033[32m"
	reset := "\033[0m"
	titulo := verde + "[Ruta Optimizada de Pallet]" + reset

	fmt.Println("\n" + titulo)
	fmt.Println("\nAquí iría la implementación para mostrar la ruta optimizada.")

	fmt.Println("\nAnalizando rutas posibles...")
	time.Sleep(1 * time.Second)
	fmt.Println("Calculando distancias...")
	time.Sleep(1 * time.Second)
	fmt.Println("Optimización completada. La mejor ruta es: A → C → B → D")

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}

func CorregirXY() {
	fmt.Print("\033[H\033[2J")

	verde := "\033[32m"
	reset := "\033[0m"
	titulo := verde + "[Corrección de X&Y]" + reset

	fmt.Println("\n" + titulo)
	fmt.Println("\nAquí iría la implementación de la corrección de coordenadas.")

	fmt.Println("\nProcesando coordenadas...")
	time.Sleep(1 * time.Second)
	fmt.Println("Conectando a Google Places API...")
	time.Sleep(1 * time.Second)
	fmt.Println("Coordenadas corregidas correctamente.")

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}

func MostrarAyuda() {
	fmt.Print("\033[H\033[2J")

	verde := "\033[32m"
	reset := "\033[0m"
	titulo := verde + "[Ayuda]" + reset

	fmt.Println("\n" + titulo)
	fmt.Println("\nEsta aplicación CLI permite realizar diversas tareas relacionadas con la gestión de coordenadas y rutas.")
	fmt.Println("\nOpciones disponibles:")
	fmt.Println("- Corregir X&Y: Herramienta para ajustar coordenadas usando Google Places")
	fmt.Println("- Mostrar ruta optimizada: Calcula la mejor ruta para un pallet")
	fmt.Println("- Obtener coordenadas: Extrae coordenadas de un pallet y las guarda en un archivo")
	fmt.Println("- Generar mapa HTML: Crea un mapa interactivo a partir de un archivo de coordenadas")
	fmt.Println("- Ayuda: Muestra esta información")

	fmt.Println("\nInstrucciones de uso:")
	fmt.Println("1. Usa las flechas ↑/↓ para navegar por el menú")
	fmt.Println("2. Presiona Enter para seleccionar una opción")
	fmt.Println("3. En cualquier momento puedes presionar q para salir")

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}
