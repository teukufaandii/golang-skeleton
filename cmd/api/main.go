package main

import (
	"golang-skeleton/config"
	handler "golang-skeleton/internal/handler/http"
	"golang-skeleton/internal/repository"
	"golang-skeleton/internal/routes"
	services "golang-skeleton/internal/service"
	database "golang-skeleton/pkg/database/postgre"
	redis "golang-skeleton/pkg/database/redis"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// config init
	cfg := config.Load()

	db, err := database.NewPostgreConn(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	database.AutoMigrate(db)
	database.SeedDatabase(db)
	database.CreateIndexes(db)

	// redis connection init
	redisClient, err := redis.NewRedisConn(cfg.RedisURL, cfg.RedisPassword, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// repo connection
	userRepo := repository.NewUserRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)

	// service connection
	authService := services.NewAuthService(*userRepo, redisRepo, cfg)

	// handler connection
	authHandler := handler.NewAuthHandler(authService)

	// routing
	router := routes.SetupRouter(
		authHandler,
		cfg.JWTSecret,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
