package handlers

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/music-shares/api/internal/models"
	"github.com/music-shares/api/internal/services"
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
        c.JSON(http.StatusBadRequest, gin.H{"error": "fichier audio requis"})
        return
    }

    // Upload vers MinIO
    objectName, err := h.storageService.UploadFile(file)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("erreur upload: %v", err)})
        return
    }

    // Créer l'entrée dans la BD avec l'ObjectName
    music := models.Music{
        ID: uuid.New().String(),
        Title: c.PostForm("title"),
        AudioFile: objectName,
        UserID: userID.(string),
    }

    // Correction ici : ajouter userID comme deuxième argument
    if err := h.musicService.CreateMusic(&music, userID.(string)); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur création musique"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"music": music})
}

// Pour le streaming
func (h *MusicHandler) Stream(c *gin.Context) {
    musicID := c.Param("id")
    
    // Récupérer la musique de la BD
    music, err := h.musicService.GetMusic(musicID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "musique non trouvée"})
        return
    }

    // Générer une URL présignée valide 1 heure
    url, err := h.storageService.GetPresignedURL(music.AudioFile, time.Hour)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "erreur génération URL"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"stream_url": url})
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

func (h *MusicHandler) ListAllFiles(c *gin.Context) {
	log.Printf("Démarrage de ListAllFiles") 
 
	files, err := h.storageService.ListFiles()
	if err != nil {
		log.Printf("Erreur lors du listage des fichiers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erreur listage fichiers: %v", err),
			"details": "Échec de la récupération des fichiers depuis MinIO", 
		})
		return
	}
 
	var response []gin.H
	for _, file := range files {
		log.Printf("Traitement du fichier: %s", file.Key)
		
		url, err := h.storageService.GetPresignedURL(file.Key, time.Hour) // Utilise GetPresignedURL à la place
		if err != nil {
			log.Printf("Erreur génération URL pour %s: %v", file.Key, err)
			continue
		}
 
		response = append(response, gin.H{
			"name": file.Key,
			"size": file.Size,
			"lastModified": file.LastModified,
			"url": url,
		})
	}
 
	log.Printf("Total fichiers traités: %d", len(response))
 
	c.JSON(http.StatusOK, gin.H{
		"files": response,
		"count": len(response),
	})
 }