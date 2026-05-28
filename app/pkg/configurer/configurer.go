package configurer

import (
	"log"
	"os"

	"github.com/Chocobo11218/go-auth-jwt/app/internal/config"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/viper"
)

func LoadConfig() *config.Config {
	LoadEnvFile()

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot read config: %v", err)
	}

	var cfg config.Config

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("cannot unmarshal config: %v", err)
	}

	return &cfg
}

func LoadEnvFile() {
	_ = godotenv.Load()
	appEnv := os.Getenv("APP_ENV")

	if appEnv != "local" {
		log.Println("skip loading secret.env")
		return
	}

	if err := godotenv.Load("./config/secret.env"); err != nil {
		log.Fatalf("cannot load secret.env: %v", err)
	}

	log.Println("secret.env loaded")
}

func LoadSecret(secret *config.SecretConfig) {
	if err := envconfig.Process("", secret); err != nil {
		log.Fatalf("cannot load secret env: %v", err)
	}
}
