package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func LoadEnv() {
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

func GetAPICredentials() (string, string) {
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

	return apiUser, apiPassword
}
