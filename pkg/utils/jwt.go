package utils

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/music-shares/api/internal/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

const (
	TokenExpiration = 24 * time.Hour
	JWTSecret       = "your-256-bit-secret" // À changer en production !
)

func GenerateToken(user *models.User) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JWTSecret))
}

func ExtractToken(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		return "", errors.New("no authorization header")
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	return parts[1], nil
}

func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GetMusicsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var musics []models.Music
		if err := db.Find(&musics).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur de récupération des musiques"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"musics": musics})
	}
}

func CreateMusicHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("user_id")

		var music models.Music
		if err := c.ShouldBindJSON(&music); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		music.UserID = userID.(string)
		if err := db.Create(&music).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create music"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"music": music})
	}
}
