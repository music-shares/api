package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors" // Changez l'import
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/music-shares/api/internal/config"
	"github.com/music-shares/api/internal/database"
	"github.com/music-shares/api/internal/handlers"
	"github.com/music-shares/api/internal/middleware"
	"github.com/music-shares/api/internal/services"
)

func main() {
    err := godotenv.Load(".env")
    if err != nil {
        log.Printf("Warning: .env file not found, using defaults")
    }
    
    cfg := config.LoadConfig()
    
    // Log de démarrage
    log.Printf("Starting server on port %s", cfg.Server.Port)
    log.Printf("Database connection: host=%s port=%s dbname=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Name)
    

	// Initialiser la DB
	db, err := database.InitDB()
	if err != nil {
		panic(err)
	}

	// Initialiser les services
	authService := services.NewAuthService(db)
	musicService := services.NewMusicService(db)
	storageService, err := services.NewStorageService()
	if err != nil {
		panic(err)
	}
	// Initialiser les handlers
	authHandler := handlers.NewAuthHandler(authService)


	// Initialiser les handlers
	musicHandler := handlers.NewMusicHandler(musicService, storageService)

	// Configuration Gin
	r := gin.Default()
    // Configuration CORS
    r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://music.okloud-hub.com",
			"https://api.okloud-hub.com",
			"http://localhost:5173",
		},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:          12 * time.Hour,
    }))
	// Routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/check-auth", middleware.Auth(), authHandler.CheckAuth) // Nouvelle route

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
			musics.GET("/:id/stream", musicHandler.Stream)
			musics.GET("/files", musicHandler.ListAllFiles)
		}
	}
    if err := r.Run(":" + cfg.Server.Port); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
