package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	LoadDotEnv()
}

func LoadDotEnv() {
	log.Println("Loading .env file")

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	requiredEnvVars := []string{"MONGODB_URI"}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Environment variable %s must not be empty. See\n\t https://github.com/konstfish/micro-mailer/README.md", envVar)
		}
	}
}

func GetEnv(variable string) string {
	return os.Getenv(variable)
}
