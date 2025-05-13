package models

type CoordInfo struct {
	Lat             float64
	Lon             float64
	VehicleLocation int
	Index           int
}

type DeliveryOrderResponse struct {
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

type Coordenada struct {
	Lat   float64
	Lon   float64
	Index int
}

type MapData struct {
	Coordenadas []Coordenada
	CentroLat   float64
	CentroLon   float64
}
