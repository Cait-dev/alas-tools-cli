package handlers

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Cait-dev/alas-tools-cli/internal/models"
)

func GenerarMapaHTML(coordenadasTXT string) {
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

	coordenadas, err := parseCoordinatesFile(coordenadasTXT)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	if len(coordenadas) == 0 {
		fmt.Println(verde + "\n[ERROR]" + reset + " No se pudieron extraer coordenadas válidas del archivo.")
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	// Calcular el centro del mapa
	var sumLat, sumLon float64
	for _, coord := range coordenadas {
		sumLat += coord.Lat
		sumLon += coord.Lon
	}
	centroLat := sumLat / float64(len(coordenadas))
	centroLon := sumLon / float64(len(coordenadas))

	datos := models.MapData{
		Coordenadas: coordenadas,
		CentroLat:   centroLat,
		CentroLon:   centroLon,
	}

	// Generar HTML
	nombreBase := strings.TrimSuffix(coordenadasTXT, ".txt")
	nombreHTML := nombreBase + ".html"

	err = generateHTMLMap(nombreHTML, datos)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " " + err.Error())
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

func parseCoordinatesFile(filePath string) ([]models.Coordenada, error) {
	contenido, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo: %w", err)
	}

	coordStr := string(contenido)
	coordStr = strings.TrimSpace(coordStr)

	coordStr = strings.TrimPrefix(coordStr, "[")
	coordStr = strings.TrimSuffix(coordStr, "]")

	paresCoordenadas := strings.Split(coordStr, "), (")

	for i := range paresCoordenadas {
		paresCoordenadas[i] = strings.Trim(paresCoordenadas[i], "()")
	}

	var coordenadas []models.Coordenada

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

		coordenadas = append(coordenadas, models.Coordenada{
			Lat:   lat,
			Lon:   lon,
			Index: i + 1,
		})
	}

	return coordenadas, nil
}

func generateHTMLMap(fileName string, datos models.MapData) error {
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
		return fmt.Errorf("error al procesar la plantilla: %w", err)
	}

	archivoHTML, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error al crear el archivo HTML: %w", err)
	}
	defer archivoHTML.Close()

	err = tmpl.Execute(archivoHTML, datos)
	if err != nil {
		return fmt.Errorf("error al generar el HTML: %w", err)
	}

	return nil
}
