package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ims-intro/pkg/common/app"
	"ims-intro/pkg/common/postgresql"
	"ims-intro/pkg/controller"
	"ims-intro/pkg/repository"
	"ims-intro/pkg/service"
	"log"
	"os"
	"path/filepath"
)

func main() {
	envPath := filepath.Join("..", ".env")
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	configurationManager := app.NewConfigurationManager()
	dbPool := postgresql.GetConnectionPool(ctx, configurationManager.PostgresqlConfig)

	userRepository := repository.NewUserRepository(dbPool)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
}
