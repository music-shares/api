// internal/handlers/auth_handler.go
package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/music-shares/api/internal/models"
	"github.com/music-shares/api/internal/services"
	"net/http"
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
	h.authService.Register(c)
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

}
