// internal/services/auth_service.go
package services

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/music-shares/api/internal/models"
	"github.com/music-shares/api/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	Db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		Db: db,
	}
}

func (s *AuthService) Register(c *gin.Context) {
    var user models.User
    
    // Log la requête reçue
    log.Printf("Tentative d'enregistrement - Données reçues : %+v", c.Request.Body)

    if err := c.ShouldBindJSON(&user); err != nil {
        log.Printf("Erreur de binding JSON: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Log après le binding
    log.Printf("Données après binding : %+v", user)

    // Génération d'un ID si non fourni
    if user.ID == "" {
        user.ID = uuid.New().String()
    }

    // Hash du mot de passe
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Erreur de hashage du mot de passe: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
        return
    }
    user.Password = string(hashedPassword)

    // Création de l'utilisateur dans la BD
    if err := s.Db.Create(&user).Error; err != nil {
        log.Printf("Erreur lors de la création en BD: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    log.Printf("Utilisateur créé avec succès: %s", user.ID)

    c.JSON(http.StatusCreated, gin.H{
        "user": user.ID,
        "message": "User created successfully",
    })
}

func (s *AuthService) Login(c *gin.Context) {
	var loginRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := s.Db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
