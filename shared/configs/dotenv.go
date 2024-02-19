package configs

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

var genericConfig = []string{
	"MONGODB_URI",
	"REDIS_URI",
}

func LoadConfig(expectedVars ...string) {
	godotenv.Load()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Error reading config file, %s", err)
	}

	expectedVars = append(expectedVars, genericConfig...)

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
