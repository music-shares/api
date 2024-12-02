package handlers

import (
   "github.com/gin-gonic/gin"
   "github.com/music-shares/api/internal/models"
   "github.com/music-shares/api/internal/services"
   "github.com/google/uuid"
   "net/http"
   "path/filepath"
)

type MusicHandler struct {
   musicService *services.MusicService
   storageService *services.StorageService
}

func NewMusicHandler(musicService *services.MusicService, storageService *services.StorageService) *MusicHandler {
   return &MusicHandler{
       musicService: musicService,
       storageService: storageService,
   }
}

func (h *MusicHandler) Create(c *gin.Context) {
   userID, exists := c.Get("user_id")
   if !exists {
       c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
       return
   }

   file, err := c.FormFile("audio")
   if err != nil {
       c.JSON(http.StatusBadRequest, gin.H{"error": "audio file required"})
       return
   }

   // Get form data
   title := c.PostForm("title")

   // Check file type
   if !isAllowedAudioFile(file.Filename) {
       c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file format"})
       return
   }

   // Upload to MinIO
   objectName, err := h.storageService.UploadFile(file)
   if err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file"})
       return
   }

   // Create music entry
   music := models.Music{
       ID: uuid.New().String(),
       Title: title,
       AudioFile: objectName,
       UserID: userID.(string),
   }

   if err := h.musicService.CreateMusic(&music, userID.(string)); err != nil {
       c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
       return
   }

   c.JSON(http.StatusCreated, gin.H{"music": music})
}

func (h *MusicHandler) GetAll(c *gin.Context) {
   musics, err := h.musicService.GetAllMusic()
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

   c.JSON(http.StatusOK, gin.H{"message": "music deleted"})
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

func isAllowedAudioFile(filename string) bool {
   allowedExt := map[string]bool{
       ".mp3": true,
       ".wav": true,
       ".ogg": true,
       ".m4a": true,
       ".flac": true,
   }
   ext := filepath.Ext(filename)
   return allowedExt[ext]
}