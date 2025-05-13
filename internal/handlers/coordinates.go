package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/Cait-dev/alas-tools-cli/internal/api"
	"github.com/Cait-dev/alas-tools-cli/internal/config"
	"github.com/Cait-dev/alas-tools-cli/internal/models"
)

func ObtenerCoordenadas() {
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

	fmt.Printf("\nConsultando API para %d pallet(s): %s...\n", len(validPalletCodes), strings.Join(validPalletCodes, ", "))

	apiUser, apiPassword := config.GetAPICredentials()
	client := api.NewClient(apiUser, apiPassword)

	sourceFields := []string{"vehicle_location", "destination.geo_location"}

	responseBody, err := client.SearchDeliveryOrders(validPalletCodes, 0, 1, sourceFields)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	var initialResponseData struct {
		Total int `json:"total"`
	}

	err = json.Unmarshal(responseBody, &initialResponseData)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al procesar la respuesta: " + err.Error())
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

	responseBody, err = client.SearchDeliveryOrders(validPalletCodes, 0, totalItems, sourceFields)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	var responseData models.DeliveryOrderResponse
	err = json.Unmarshal(responseBody, &responseData)
	if err != nil {
		fmt.Println(verde + "\n[ERROR]" + reset + " Error al procesar la respuesta: " + err.Error())
		fmt.Println("\nPresiona Enter para volver al menú principal...")
		fmt.Scanln()
		return
	}

	var coordInfos []models.CoordInfo
	for i, item := range responseData.Items {
		lat := item.Destination.GeoLocation.Lat
		lon := item.Destination.GeoLocation.Lon
		vehicleLoc := item.VehicleLocation

		if lat != 0 && lon != 0 {
			coordInfos = append(coordInfos, models.CoordInfo{
				Lat:             lat,
				Lon:             lon,
				VehicleLocation: vehicleLoc,
				Index:           i,
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

	processAndSaveCoordinates(coordInfos, validPalletCodes)
}

func processAndSaveCoordinates(coordInfos []models.CoordInfo, palletCodes []string) {
	verde := "\033[32m"
	reset := "\033[0m"
	reader := bufio.NewReader(os.Stdin)

	var coordinates []string
	for i, info := range coordInfos {
		coordinates = append(coordinates, fmt.Sprintf("(%.7f, %.7f) /* Orden #%d, Vehicle Location: %d */",
			info.Lat, info.Lon, i+1, info.VehicleLocation))
	}

	var coordinatesClean []string
	for _, info := range coordInfos {
		coordinatesClean = append(coordinatesClean, fmt.Sprintf("(%.7f, %.7f)", info.Lat, info.Lon))
	}

	coordinatesStr := "[" + strings.Join(coordinates, ", ") + "]"
	coordinatesCleanStr := "[" + strings.Join(coordinatesClean, ", ") + "]"

	var filename string
	if len(palletCodes) == 1 {
		filename = fmt.Sprintf("coordenadas_%s.txt", palletCodes[0])
	} else {
		filename = fmt.Sprintf("coordenadas_multiple_%d_pallets.txt", len(palletCodes))
	}

	err := ioutil.WriteFile(filename, []byte(coordinatesStr), 0644)
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
		GenerarMapaHTML(filenameClean)
	}

	fmt.Println("\nPresiona Enter para volver al menú principal...")
	fmt.Scanln()
}
