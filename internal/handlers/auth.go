// internal/handlers/auth_handler.go
package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/music-shares/api/internal/models"
	"github.com/music-shares/api/internal/services"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
    var user models.User
    
    // Log de debug
    fmt.Printf("Début de l'enregistrement\n")

    if err := c.ShouldBindJSON(&user); err != nil {
        fmt.Printf("Erreur de binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Log des données reçues
    fmt.Printf("Données reçues: %+v\n", user)

    // Hashage du mot de passe
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Printf("Erreur de hashage: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
        return
    }

    // Préparation de l'utilisateur
    user.Password = string(hashedPassword)
    user.ID = uuid.New().String()
    user.CreatedAt = time.Now()
    user.UpdatedAt = time.Now()

    // Création dans la base de données
    if err := h.authService.Db.Create(&user).Error; err != nil {
        fmt.Printf("Erreur de création en DB: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

    // Succès
    c.JSON(http.StatusCreated, gin.H{
        "message": "User created successfully",
        "user": gin.H{
            "id": user.ID,
            "email": user.Email,
            "username": user.Username,
        },
    })
}
func (h *AuthHandler) Login(c *gin.Context) {
	h.authService.Login(c)
}

func (h *AuthHandler) GetAllUsers(c *gin.Context) {
	var users []models.User
	if err := h.authService.Db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération des utilisateurs"})
		return
	}
	// Debug: afficher les utilisateurs dans les logs
	fmt.Printf("Nombre d'utilisateurs trouvés: %d\n", len(users))
	for _, user := range users {
		fmt.Printf("User: ID=%s, Email=%s\n", user.ID, user.Email)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (h *AuthHandler) Logout(c *gin.Context) {
    // Implémentation du logout si nécessaire
    c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}

func (h *AuthHandler) CheckAuth(c *gin.Context) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var user models.User
    if err := h.authService.Db.First(&user, "id = ?", userID).Error; err != nil { // Changé db en DB
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "user": gin.H{
            "id": user.ID,
            "email": user.Email,
            "username": user.Username,
        },
    })
}