package commons

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Aviso: Arquivo .env não encontrado. Verificando variáveis de ambiente do sistema...")
	}
}

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}

	fmt.Printf("Aviso: A variável %s não está definida.\n", key)
	return ""
}
