package main

import (
	"bytes"
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error al leer el archivo .env: %v\n", err)
	}
}

var (
	version = "dev"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Foreground(lipgloss.Color("#00FF00")).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Foreground(lipgloss.Color("#FFFFFF"))
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#00FF00")).Bold(true)
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4).Foreground(lipgloss.Color("#00FF00"))
)

type menuItem struct {
	title, desc string
}

func (m menuItem) Title() string       { return m.title }
func (m menuItem) Description() string { return m.desc }
func (m menuItem) FilterValue() string { return m.title }

type model struct {
	list     list.Model
	quitting bool
	action   int
}

func initialModel() model {
	loadEnv()
	items := []list.Item{
		menuItem{title: "Opción 1: Corregir X&Y", desc: "Herramienta para corregir coordenadas usando Google Places"},
		menuItem{title: "Opción 2: Mostrar una ruta optimizada de un pallet", desc: "Compara las rutas"},
		menuItem{title: "Opción 3: Obtener coordenadas", desc: "Extrae coordenadas de un pallet y las guarda en un archivo"},
		menuItem{title: "Opción 4: Generar mapa HTML", desc: "Crea un mapa interactivo a partir de un archivo de coordenadas"},
		menuItem{title: "Opción 5: Ayuda", desc: "Muestra la información de ayuda"},
		menuItem{title: "Salir", desc: "Salir de la aplicación"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = selectedItemStyle
	delegate.Styles.NormalTitle = itemStyle

	l := list.New(items, delegate, 80, 20)
	l.Title = "Menú Principal"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	return model{
		list:     l,
		quitting: false,
		action:   -1,
	}
}

func mostrarPantallaInicio() {
	fmt.Print("\033[H\033[2J")
	asciiArt := `
/$$$$$$  /$$                         /$$$$$$$$                  /$$              /$$$$$$  /$$ /$$
/$$__  $$| $$                        |__  $$__/                 | $$             /$$__  $$| $$|__/
| $$  \ $$| $$ /$$   /$$ /$$$$$$$       | $$  /$$$$$$   /$$$$$$ | $$  /$$$$$$$ | $$  \__/| $$ /$$
| $$$$$$$$| $$|  $$ /$$/| $$__  $$      | $$ /$$__  $$ /$$__  $$| $$ /$$_____/ | $$      | $$| $$
| $$__  $$| $$ \  $$$$/ | $$  \ $$      | $$| $$  \ $$| $$  \ $$| $$|  $$$$$$ | $$      | $$| $$
| $$  | $$| $$  >$$  $$ | $$  | $$      | $$| $$  | $$| $$  | $$| $$ \____  $$| $$    $$| $$| $$
| $$  | $$| $$ /$$/\  $$| $$  | $$      | $$|  $$$$$$/|  $$$$$$/| $$ /$$$$$$$/|  $$$$$$/| $$| $$
|__/  |__/|__/|__/  \__/|__/  |__/      |__/ \______/  \______/ |__/|_______/  \______/ |__/|__/
`
	subtitulo := `
┌─────────────────────────────────┐
│ Entertainment - Development     │
│ Productivity - & much more!     │
└─────────────────────────────────┘
`
	verde := "\033[32m"
	reset := "\033[0m"
	fmt.Println(verde + asciiArt + subtitulo + reset)

	fmt.Println("\nBienvenido a Alas-Tools-Cli v1.1.1")
	fmt.Println("─────────────────────────────")
	fmt.Println("Use the arrow keys to navigate: ↑ ↓")
	time.Sleep(2 * time.Second)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				if strings.Contains(i.title, "Salir") {
					m.quitting = true
					return m, tea.Quit
				}

				m.action = m.list.Index()
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func (m model) View() string {
	if m.quitting {
		return quitTextStyle.Render("¡Gracias por usar Alas-Tools-Cli!")
	}

	return docStyle.Render(m.list.View())
}

func main() {

	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("Alas-Tools-Cli versión %s\n", version)
		return
	}

	mostrarPantallaInicio()

	salir := false
	for !salir {
		m := initialModel()
		p := tea.NewProgram(m)
		finalModel, err := p.Run()
		if err != nil {
			fmt.Printf("Error al iniciar la aplicación: %v\n", err)
			os.Exit(1)
		}

		if finalModel, ok := finalModel.(model); ok {
			if finalModel.quitting {
				salir = true
			} else {
				switch finalModel.action {
				case 0:
					corregirXY()
				case 1:
					mostrarRutaOptimizada()
				case 2:
					obtenerCoordenadas()
				case 3:
					generarMapaHTML("")
				case 4:
					mostrarAyuda()
				}

				fmt.Print("\033[H\033[2J")
				fmt.Println("\nVolviendo al menú principal...")
				time.Sleep(1 * time.Second)
			}
		}
	}

	fmt.Println("\n¡Hasta pronto! Gracias por usar Alas-Tools-Cli.")
}

func corregirXY() {
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

func mostrarRutaOptimizada() {
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

func mostrarAyuda() {
	fmt.Print("\033[H\033[2J")

	verde := "\033[32m"
	reset := "\033[0m"
	titulo := verde + "[Ayuda]" + reset

	fmt.Println("\n" + titulo)
	fmt.Println("\nEsta aplicación CLI permite realizar diversas tareas relacionadas con la gestión de coordenadas y rutas.")
	fmt.Println("\nOpciones disponibles:")
	fmt.Println("- Corregir X&Y: Herramienta para ajustar coordenadas usando Google Places")
	fmt.Println("- Mostrar ruta optimizada: Calcula la mejor ruta para un pallet")
	fmt.Println("- Ayuda: Muestra esta información")

	fmt.Println("\nInstrucciones de uso:")
	fmt.Println("1. Usa las flechas ↑/↓ para navegar por el menú")
	fmt.Println("2. Presiona Enter para seleccionar una opción")
	fmt.Println("3. En cualquier momento puedes presionar q para salir")

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}

func obtenerCoordenadas() {
	fmt.Print("\033[H\033[2J")

	verde := "\033[32m"
	reset := "\033[0m"
	titulo := verde + "[Obtener Coordenadas]" + reset

	fmt.Println("\n" + titulo)
	fmt.Println("\nEsta herramienta extrae coordenadas de las órdenes asociadas a un pallet.")

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("\nIngrese el código de pallet (ej. pl202505danl001) o varios separados por comas: ")
	palletInput, _ := reader.ReadString('\n')
	palletInput = strings.TrimSpace(palletInput)

	palletCodes := strings.Split(palletInput, ",")
	for i, code := range palletCodes {
		palletCodes[i] = strings.TrimSpace(code)
	}

	var validPalletCodes []string
	for _, code := range palletCodes {
		if code != "" {
			validPalletCodes = append(validPalletCodes, code)
		}
	}

	if len(validPalletCodes) == 0 {
		fmt.Println(verde + "\n[ERROR]" + reset + " Debe ingresar al menos un código de pallet válido.")
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	fmt.Print("\nIngrese el nombre de usuario: ")
	userName, _ := reader.ReadString('\n')
	userName = strings.TrimSpace(userName)

	if userName == "" {
		userName = "raul.sepulveda"
		fmt.Println("Se usará el usuario predeterminado: " + userName)
	}

	fmt.Printf("\nConsultando API para %d pallet(s): %s...\n", len(validPalletCodes), strings.Join(validPalletCodes, ", "))

	apiUser := os.Getenv("ALAS_API_USER")
	apiPassword := os.Getenv("ALAS_API_PASSWORD")

	if apiUser == "" {
		apiUser = "dev_user"
		fmt.Println("Advertencia: ALAS_API_USER no está configurada, usando valor predeterminado para desarrollo")
	}
	if apiPassword == "" {
		apiPassword = "dev_password"
		fmt.Println("Advertencia: ALAS_API_PASSWORD no está configurada, usando valor predeterminado para desarrollo")
	}

	initialRequestBody := struct {
		PalletCodes  []string `json:"pallet_codes"`
		PageNumber   int      `json:"page_number"`
		PageSize     int      `json:"page_size"`
		SourceFields []string `json:"source_fields"`
	}{
		PalletCodes:  validPalletCodes,
		PageNumber:   0,
		PageSize:     1,
		SourceFields: []string{"vehicle_location", "destination.geo_location"},
	}

	initialRequestJSON, err := json.Marshal(initialRequestBody)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al crear la petición inicial: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	client := &http.Client{}
	initialReq, err := http.NewRequest("POST", "https://api.alasxpress.com/delivery/delivery-orders/cl/_search", bytes.NewBuffer(initialRequestJSON))
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al crear la petición inicial: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	initialReq.Header.Set("Content-Type", "application/json")
	initialReq.SetBasicAuth(apiUser, apiPassword)

	fmt.Println("Obteniendo información de paginación...")
	initialResp, err := client.Do(initialReq)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al conectar con la API: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}
	defer initialResp.Body.Close()

	var initialResponseData struct {
		Total int `json:"total"`
	}

	initialBody, _ := ioutil.ReadAll(initialResp.Body)
	if initialResp.StatusCode != 200 {
		fmt.Printf("%s\n[ERROR]%s Código de estado: %d - %s\n", verde, reset, initialResp.StatusCode, string(initialBody))
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	err = json.Unmarshal(initialBody, &initialResponseData)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al procesar la respuesta inicial: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	totalItems := initialResponseData.Total
	if totalItems == 0 {
		fmt.Println(verde + "\n[AVISO]" + reset + " No se encontraron órdenes para los pallets proporcionados.")
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	fmt.Printf("Se encontraron un total de %d órdenes. Obteniendo coordenadas...\n", totalItems)

	mainRequestBody := struct {
		PalletCodes  []string `json:"pallet_codes"`
		PageNumber   int      `json:"page_number"`
		PageSize     int      `json:"page_size"`
		SourceFields []string `json:"source_fields"`
	}{
		PalletCodes:  validPalletCodes,
		PageNumber:   0,
		PageSize:     totalItems,
		SourceFields: []string{"vehicle_location", "destination.geo_location"},
	}

	mainRequestJSON, err := json.Marshal(mainRequestBody)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al crear la petición principal: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	mainReq, err := http.NewRequest("POST", "https://api.alasxpress.com/delivery/delivery-orders/cl/_search", bytes.NewBuffer(mainRequestJSON))
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al crear la petición principal: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	mainReq.Header.Set("Content-Type", "application/json")
	mainReq.SetBasicAuth(apiUser, apiPassword)

	fmt.Println("Obteniendo coordenadas detalladas...")
	mainResp, err := client.Do(mainReq)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al conectar con la API: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}
	defer mainResp.Body.Close()

	mainBody, err := ioutil.ReadAll(mainResp.Body)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al leer la respuesta: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	if mainResp.StatusCode != 200 {
		fmt.Printf("%s\n[ERROR]%s Código de estado: %d - %s\n", verde, reset, mainResp.StatusCode, string(mainBody))
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	var responseData struct {
		Total int `json:"total"`
		Items []struct {
			Destination struct {
				GeoLocation struct {
					Lat float64 `json:"lat"`
					Lon float64 `json:"lon"`
				} `json:"geo_location"`
			} `json:"destination"`
			VehicleLocation int `json:"vehicle_location"`
		} `json:"items"`
	}

	err = json.Unmarshal(mainBody, &responseData)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al procesar la respuesta: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	type CoordInfo struct {
		Lat            float64
		Lon            float64
		VehicleLocation int
		Index          int
	}

	var coordInfos []CoordInfo
	for i, item := range responseData.Items {
		lat := item.Destination.GeoLocation.Lat
		lon := item.Destination.GeoLocation.Lon
		vehicleLoc := item.VehicleLocation

		if lat != 0 && lon != 0 {
			coordInfos = append(coordInfos, CoordInfo{
				Lat:            lat,
				Lon:            lon,
				VehicleLocation: vehicleLoc,
				Index:          i,
			})
		}
	}

	if len(coordInfos) == 0 {
		fmt.Println(verde + "\n[AVISO]" + reset + " No se encontraron coordenadas válidas para los pallets proporcionados.")
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	sort.Slice(coordInfos, func(i, j int) bool {
		return coordInfos[i].VehicleLocation < coordInfos[j].VehicleLocation
	})

	var coordinates []string
	for i, info := range coordInfos {
		coordinates = append(coordinates, fmt.Sprintf("(%.7f, %.7f) /* Orden #%d, Vehicle Location: %d */", 
			info.Lat, info.Lon, i + 1, info.VehicleLocation))
	}

	var coordinatesClean []string
	for _, info := range coordInfos {
		coordinatesClean = append(coordinatesClean, fmt.Sprintf("(%.7f, %.7f)", info.Lat, info.Lon))
	}

	coordinatesStr := "[" + strings.Join(coordinates, ", ") + "]"
	coordinatesCleanStr := "[" + strings.Join(coordinatesClean, ", ") + "]"

	var filename string
	if len(validPalletCodes) == 1 {
		filename = fmt.Sprintf("coordenadas_%s.txt", validPalletCodes[0])
	} else {
		filename = fmt.Sprintf("coordenadas_multiple_%d_pallets.txt", len(validPalletCodes))
	}

	err = ioutil.WriteFile(filename, []byte(coordinatesStr), 0644)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al escribir el archivo: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	filenameClean := strings.TrimSuffix(filename, ".txt") + "_clean.txt"
	err = ioutil.WriteFile(filenameClean, []byte(coordinatesCleanStr), 0644)
	if err != nil {
		fmt.Println(verde + "\n[AVISO]" + reset + " Error al escribir el archivo limpio: " + err.Error())
	}

	fmt.Printf("\n%s[ÉXITO]%s Se encontraron %d coordenadas ordenadas por Vehicle Location.\n", verde, reset, len(coordinates))
	fmt.Printf("Se ha creado el archivo %s con las coordenadas en el formato solicitado.\n", filename)
	fmt.Printf("También se creó %s con un formato compatible para otras herramientas.\n", filenameClean)

	fmt.Print("\n¿Desea generar un mapa HTML con estas coordenadas? (s/n): ")
	respuesta, _ := reader.ReadString('\n')
	respuesta = strings.TrimSpace(respuesta)

	if strings.ToLower(respuesta) == "s" || strings.ToLower(respuesta) == "si" {
		generarMapaHTML(filenameClean)
	}

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}

func generarMapaHTML(coordenadasTXT string) {
	fmt.Print("\033[H\033[2J")

	verde := "\033[32m"
	reset := "\033[0m"
	titulo := verde + "[Generar Mapa HTML]" + reset

	fmt.Println("\n" + titulo)
	fmt.Println("\nEsta herramienta genera un mapa HTML interactivo a partir de un archivo de coordenadas.")

	if coordenadasTXT == "" {
		fmt.Print("\nIngrese la ruta del archivo de coordenadas (ej. coordenadas_pl202505danl001.txt): ")
		fmt.Scanln(&coordenadasTXT)
	}

	if coordenadasTXT == "" {
		fmt.Println(verde + "\n[ERROR]" + reset + " Debe proporcionar un archivo de coordenadas.")
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	if _, err := os.Stat(coordenadasTXT); os.IsNotExist(err) {
		fmt.Println(verde + "\n[ERROR]" + reset + " El archivo especificado no existe: " + coordenadasTXT)
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	contenido, err := ioutil.ReadFile(coordenadasTXT)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al leer el archivo: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	coordStr := string(contenido)
	coordStr = strings.TrimSpace(coordStr)

	coordStr = strings.TrimPrefix(coordStr, "[")
	coordStr = strings.TrimSuffix(coordStr, "]")

	paresCoordenadas := strings.Split(coordStr, "), (")

	for i := range paresCoordenadas {
		paresCoordenadas[i] = strings.Trim(paresCoordenadas[i], "()")
	}

	type Coordenada struct {
		Lat   float64
		Lon   float64
		Index int
	}
	var coordenadas []Coordenada

	for i, par := range paresCoordenadas {
		partes := strings.Split(par, ", ")
		if len(partes) != 2 {
			continue
		}

		lat, err1 := strconv.ParseFloat(partes[0], 64)
		lon, err2 := strconv.ParseFloat(partes[1], 64)

		if err1 != nil || err2 != nil {
			continue
		}

		coordenadas = append(coordenadas, Coordenada{
			Lat:   lat,
			Lon:   lon,
			Index: i + 1,
		})
	}

	if len(coordenadas) == 0 {
		fmt.Println(verde + "\n[ERROR]" + reset + " No se pudieron extraer coordenadas válidas del archivo.")
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	var sumLat, sumLon float64
	for _, coord := range coordenadas {
		sumLat += coord.Lat
		sumLon += coord.Lon
	}
	centroLat := sumLat / float64(len(coordenadas))
	centroLon := sumLon / float64(len(coordenadas))

	htmlTemplate := `<!DOCTYPE html>
<html lang="es">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Mapa de Coordenadas</title>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.css" />
    <style>
        body {
            margin: 0;
            padding: 0;
            font-family: Arial, sans-serif;
        }
        #map {
            height: 600px;
            width: 100%;
        }
        .info-panel {
            padding: 10px;
            background: white;
            border-radius: 5px;
            box-shadow: 0 0 15px rgba(0,0,0,0.2);
            margin-bottom: 10px;
        }
        .button-container {
            padding: 10px;
            text-align: center;
        }
        button {
            padding: 10px 15px;
            background-color: #4CAF50;
            color: white;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            margin: 5px;
        }
        button:hover {
            background-color: #45a049;
        }
    </style>
</head>
<body>
    <div class="info-panel">
        <h1>Mapa de Coordenadas</h1>
    </div>
    
    <div class="button-container">
        <button id="toggle-line">Mostrar/Ocultar Línea</button>
        <button id="fit-bounds">Ajustar Vista</button>
    </div>
    
    <div id="map"></div>
    
    <script src="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.js"></script>
    
    <script>
        // Inicializar el mapa
        const map = L.map('map').setView([{{.CentroLat}}, {{.CentroLon}}], 13);
        
        // Añadir capa principal de OpenStreetMap
        const baseMap = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
            opacity: 1
        }).addTo(map);
        
        // Añadir capa humanitaria (con más detalles urbanos)
        const hotMap = L.tileLayer('https://{s}.tile.openstreetmap.fr/hot/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors, Tiles style by <a href="https://www.hotosm.org/" target="_blank">Humanitarian OpenStreetMap Team</a>',
            opacity: 0.7
        }).addTo(map);
        
        // Crear un grupo para todos los marcadores
        const markersGroup = L.layerGroup().addTo(map);
        
        // Coordenadas
        const coordinates = [
            {{range .Coordenadas}}
            [{{.Lat}}, {{.Lon}}],
            {{end}}
        ];
        
        // Añadir marcadores con números
        coordinates.forEach((coord, index) => {
            // Crear un div personalizado con un círculo de fondo azul claro y número
            const numberIcon = L.divIcon({
                html: '<div style="background-color: #3399ff; color: white; border-radius: 50%; width: 24px; height: 24px; display: flex; align-items: center; justify-content: center; font-weight: bold; box-shadow: 0 0 3px rgba(0,0,0,0.5);">' + (index + 1) + '</div>',
                className: '',
                iconSize: [24, 24],
                iconAnchor: [12, 12]
            });
            
            // Crear el marcador y asignar el icono personalizado
            const marker = L.marker(coord, {
                icon: numberIcon
            });
            
            // Añadir popup con información
            marker.bindPopup('<b>Punto ' + (index + 1) + '</b><br>Lat: ' + coord[0] + '<br>Lon: ' + coord[1]);
            
            // Añadir el marcador al grupo
            markersGroup.addLayer(marker);
        });
        
        // Crear una línea que conecta todos los puntos
        const polyline = L.polyline(coordinates, {
            color: 'red',
            weight: 2,
            opacity: 0.9
        }).addTo(map);
        
        // Ajustar el mapa para mostrar todos los marcadores
        map.fitBounds(markersGroup.getBounds());
        
        // Función para mostrar/ocultar la línea
        let lineVisible = true;
        document.getElementById('toggle-line').addEventListener('click', function() {
            if (lineVisible) {
                map.removeLayer(polyline);
            } else {
                map.addLayer(polyline);
            }
            lineVisible = !lineVisible;
        });
        
        // Función para ajustar la vista
        document.getElementById('fit-bounds').addEventListener('click', function() {
            map.fitBounds(markersGroup.getBounds());
        });
    </script>
</body>
</html>`

	tmpl, err := template.New("mapa").Parse(htmlTemplate)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al procesar la plantilla: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	datos := struct {
		Coordenadas []Coordenada
		CentroLat   float64
		CentroLon   float64
	}{
		Coordenadas: coordenadas,
		CentroLat:   centroLat,
		CentroLon:   centroLon,
	}

	nombreBase := strings.TrimSuffix(coordenadasTXT, ".txt")
	nombreHTML := nombreBase + ".html"

	archivoHTML, err := os.Create(nombreHTML)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al crear el archivo HTML: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}
	defer archivoHTML.Close()

	err = tmpl.Execute(archivoHTML, datos)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al generar el HTML: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	fmt.Printf("\n%s[ÉXITO]%s Se ha generado el mapa con %d puntos.\n", verde, reset, len(coordenadas))
	fmt.Printf("Archivo HTML creado: %s\n", nombreHTML)
	fmt.Println("\nPuedes abrir este archivo en cualquier navegador para ver el mapa interactivo.")

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}
