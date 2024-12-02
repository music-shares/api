package database

import (
	"fmt"
	"github.com/music-shares/api/internal/config"
	"github.com/music-shares/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	// Charger la configuration qui utilise les variables d'environnement
	cfg := config.LoadConfig()

	// Utiliser GetDBDSN pour obtenir la chaîne de connexion
	dsn := cfg.GetDBDSN()

	fmt.Printf("Connecting to database with DSN: %s\n", dsn)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("erreur de connexion à la base de données: %v", err)
	}

	err = db.AutoMigrate(&models.User{}, &models.Music{})
	if err != nil {
		return nil, fmt.Errorf("erreur de migration: %v", err)
	}

	return db, nil
}
