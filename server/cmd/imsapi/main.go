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
	err = userRepository.EnsureUserSchema()
	if err != nil {
		log.Fatalf("Failed to ensure user schema: %v", err)
	}
	err = userRepository.EnsureAdminExists()
	if err != nil {
		log.Fatalf("Failed to ensure admin user: %v", err)
	}
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)

	productRepository := repository.NewProductRepository(dbPool)
	if productSchemaEnsurer, ok := productRepository.(interface{ EnsureProductSchema() error }); ok {
		err = productSchemaEnsurer.EnsureProductSchema()
		if err != nil {
			log.Fatalf("Failed to ensure product schema: %v", err)
		}
	}
	productService := service.NewProductService(productRepository)
	productController := controller.NewProductController(productService)

	cartRepository := repository.NewCartRepository(dbPool)
	err = cartRepository.EnsureCartSchema()
	if err != nil {
		log.Fatalf("Failed to ensure cart schema: %v", err)
	}
	cartService := service.NewCartService(cartRepository, productRepository)
	cartController := controller.NewCartController(cartService)

	orderRepository := repository.NewOrderRepository(dbPool)
	err = orderRepository.EnsureOrderSchema()
	if err != nil {
		log.Fatalf("Failed to ensure order schema: %v", err)
	}
	orderItemRepository := repository.NewOrderItemRepository(dbPool)
	orderService := service.NewOrderService(orderRepository, orderItemRepository, productRepository, cartRepository)
	orderController := controller.NewOrderController(orderService)

	analyticsRepository := repository.NewAnalyticsRepository(dbPool)
	err = analyticsRepository.EnsureSeedData()
	if err != nil {
		log.Fatalf("Failed to ensure analytics seed data: %v", err)
	}
	analyticsService := service.NewAnalyticsService(analyticsRepository)
	analyticsController := controller.NewAnalyticsController(analyticsService)

	userAnalyticsRepository := repository.NewUserAnalyticsRepository(dbPool)
	userAnalyticsService := service.NewUserAnalyticsService(userAnalyticsRepository)
	userAnalyticsController := controller.NewUserAnalyticsController(userAnalyticsService)

	e := echo.New()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowCredentials: true,
	}))

	userController.RegisterUserRoutes(e)
	productController.RegisterProductRoutes(e)
	cartController.RegisterCartRoutes(e)
	orderController.RegisterOrderRoutes(e)
	analyticsController.RegisterAnalyticsRoutes(e)
	userAnalyticsController.RegisterUserAnalyticsRoutes(e)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Server is running on port", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
