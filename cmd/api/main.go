package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/music-shares/api/internal/config"
	"github.com/music-shares/api/internal/database"
	"github.com/music-shares/api/internal/handlers"
	"github.com/music-shares/api/internal/middleware"
	"github.com/music-shares/api/internal/services"
	"log"
)

func main() {
	// Charger la configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	cfg := config.LoadConfig()

	// Initialiser la DB
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// Initialiser les services
	authService := services.NewAuthService(db)
	//	musicService := services.NewMusicService(db)

	// Initialiser les handlers
	authHandler := handlers.NewAuthHandler(authService)
	//	musicHandler := handlers.NewMusicHandler(musicService)

	musicService := services.NewMusicService(db)

	// Initialiser les handlers
	musicHandler := handlers.NewMusicHandler(musicService)

	// Configuration Gin
	r := gin.Default()

	// Routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	api := r.Group("/api")
	api.Use(middleware.Auth())
	{
		api.GET("/users", authHandler.GetAllUsers)
		musics := api.Group("/musics")
		{
			musics.POST("/", musicHandler.Create)
			musics.GET("/", musicHandler.GetAll)
			musics.GET("/:id", musicHandler.GetOne)
			musics.PUT("/:id", musicHandler.Update)
			musics.DELETE("/:id", musicHandler.Delete)
			musics.GET("/user", musicHandler.GetUserMusics)
		}
	}
	r.Run(":" + cfg.Server.Port)
}
