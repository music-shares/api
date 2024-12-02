package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/music-shares/api/internal/models"
	"github.com/music-shares/api/internal/services"
	"net/http"
)

type MusicHandler struct {
	musicService *services.MusicService
}

func NewMusicHandler(musicService *services.MusicService) *MusicHandler {
	return &MusicHandler{musicService: musicService}
}

func (h *MusicHandler) Create(c *gin.Context) {
	var music models.Music
	if err := c.ShouldBindJSON(&music); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := h.musicService.CreateMusic(&music, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"music": music})
}

func (h *MusicHandler) GetAll(c *gin.Context) {
	musics, err := h.musicService.GetAllMusics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"musics": musics})
}

func (h *MusicHandler) GetOne(c *gin.Context) {
	id := c.Param("id")

	music, err := h.musicService.GetMusic(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"music": music})
}

func (h *MusicHandler) Update(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	var music models.Music
	if err := c.ShouldBindJSON(&music); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.musicService.UpdateMusic(id, &music, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"music": music})
}

func (h *MusicHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	if err := h.musicService.DeleteMusic(id, userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"music": nil})
}

func (h *MusicHandler) GetUserMusics(c *gin.Context) {
	userID, _ := c.Get("user_id")

	musics, err := h.musicService.GetUserMusic(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"musics": musics})
}
