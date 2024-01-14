package configs

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func init() {
	LoadConfig()
}

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	expectedVars := []string{
		"MONGODB_URI",
		"RABBITMQ_URI",
	}

	for _, v := range expectedVars {
		log.Println(v, viper.GetString(v))
		if viper.GetString(v) == "" {
			log.Fatalf("Environment variable %s must not be empty. See\n\t https://github.com/konstfish/angler/README.md", v)
		}
	}
}

func GetConfigVar(variable string) string {
	return viper.GetString(variable)
}
