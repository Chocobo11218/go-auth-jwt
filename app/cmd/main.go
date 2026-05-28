package main

import (
	"fmt"

	"github.com/Chocobo11218/go-auth-jwt/app/pkg/configurer"
	"github.com/Chocobo11218/go-auth-jwt/app/pkg/logger"
)

func main() {
	cfg := configurer.LoadConfig()

	configurer.LoadSecret(&cfg.Secret)

	fmt.Println("App Name:", cfg.App.Name)
	fmt.Println("App Port:", cfg.App.Port)
	fmt.Println("DB Host:", cfg.DB.Host)
	fmt.Println("DB User:", cfg.DB.User)
	fmt.Println("DB Password:", cfg.Secret.DBPassword)
	fmt.Println("JWT Expire:", cfg.JWT.Expire)

	logger.Info("Configuration loaded successfully")
}
