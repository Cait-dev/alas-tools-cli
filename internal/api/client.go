package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	BaseURL  string
	Username string
	Password string
}

func NewClient(username, password string) *Client {
	return &Client{
		BaseURL:  "https://api.alasxpress.com",
		Username: username,
		Password: password,
	}
}

func (c *Client) SearchDeliveryOrders(palletCodes []string, pageNumber, pageSize int, sourceFields []string) ([]byte, error) {
	requestBody := struct {
		PalletCodes  []string `json:"pallet_codes"`
		PageNumber   int      `json:"page_number"`
		PageSize     int      `json:"page_size"`
		SourceFields []string `json:"source_fields"`
	}{
		PalletCodes:  palletCodes,
		PageNumber:   pageNumber,
		PageSize:     pageSize,
		SourceFields: sourceFields,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %w", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", c.BaseURL+"/delivery/delivery-orders/cl/_search", bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("error al crear la petición: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(c.Username, c.Password)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error al conectar con la API: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer la respuesta: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("código de estado: %d - %s", resp.StatusCode, string(body))
	}

	return body, nil
}
