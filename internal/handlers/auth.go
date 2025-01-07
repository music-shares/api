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
	"github.com/music-shares/api/pkg/utils"
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
    // Structure pour le request body
    var registerRequest struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }
    
    // Log début
    fmt.Printf("Début de l'enregistrement\n")

    // Parse le JSON
    if err := c.ShouldBindJSON(&registerRequest); err != nil {
        fmt.Printf("Erreur de binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Log des données reçues
    fmt.Printf("Données reçues - Username: %s, Email: %s, Password: %s\n", 
        registerRequest.Username, 
        registerRequest.Email, 
        registerRequest.Password)

    // Hash du mot de passe
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost)
    if err != nil {
        fmt.Printf("Erreur de hashage: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error processing password"})
        return
    }

    // Log du hash
    fmt.Printf("Hash généré: %s\n", string(hashedPassword))

    // Création de l'utilisateur
    user := models.User{
        ID:        uuid.New().String(),
        Username:  registerRequest.Username,
        Email:     registerRequest.Email,
        Password:  string(hashedPassword),
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Sauvegarde en DB
    if err := h.authService.Db.Create(&user).Error; err != nil {
        fmt.Printf("Erreur de création en DB: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
        return
    }

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
    // Structure pour recevoir les données de login
    var loginRequest struct {
        Email    string `json:"email" binding:"required,email"`
        Password string `json:"password" binding:"required"`
    }

    // Log de debug
    fmt.Printf("Tentative de connexion\n")

    // Parse du JSON
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        fmt.Printf("Erreur de binding JSON: %v\n", err)
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data"})
        return
    }

    // Log des données reçues (ne pas logger le mot de passe en production)
    fmt.Printf("Email reçu: %s\n", loginRequest.Email)

    // Recherche de l'utilisateur
    var user models.User
    if err := h.authService.Db.Where("email = ?", loginRequest.Email).First(&user).Error; err != nil {
        fmt.Printf("Utilisateur non trouvé: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Vérification du mot de passe
    fmt.Printf("Mot de passe stocké (hash): %s\n", user.Password)
    fmt.Printf("Mot de passe fourni: %s\n", loginRequest.Password)

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
        fmt.Printf("Détails de l'erreur bcrypt: %v\n", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
        return
    }

    // Génération du token JWT
    token, err := utils.GenerateToken(&user)
    if err != nil {
        fmt.Printf("Erreur de génération du token: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
        return
    }

    // Succès
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": gin.H{
            "id": user.ID,
            "email": user.Email,
            "username": user.Username,
        },
    })
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

func (h *AuthHandler) GetUserByEmail(c *gin.Context) {
    email := c.Query("email")
    var user models.User
    
    result := h.authService.Db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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