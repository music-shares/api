package database

import (
	"fmt"
	"github.com/music-shares/api/internal/config"
	"github.com/music-shares/api/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)
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
