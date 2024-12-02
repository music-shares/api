package services

import (
	"errors"
	"github.com/google/uuid"
	"github.com/music-shares/api/internal/models"
	"gorm.io/gorm"
	"time"
)

type MusicService struct {
	db *gorm.DB
}

func NewMusicService(db *gorm.DB) *MusicService {
	return &MusicService{
		db: db,
	}
}
func (s *MusicService) CreateMusic(music *models.Music, userID string) error {
	music.ID = uuid.New().String()
	music.UserID = userID
	music.CreatedAt = time.Now()
	music.UpdatedAt = time.Now()
	return s.db.Create(&music).Error
}

func (s *MusicService) GetMusic(id string) (*models.Music, error) {
	var music models.Music
	if err := s.db.Preload("User").First(&music, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("music not found")
		}
		return nil, err
	}
	return &music, nil
}

func (s *MusicService) GetAllMusic() ([]models.Music, error) {
	var musics []models.Music
	if err := s.db.Preload("User").Find(&musics).Error; err != nil {
		return nil, err
	}
	return musics, nil
}

func (s *MusicService) UpdateMusic(id string, music *models.Music, userID string) error {
	existingMusic, err := s.GetMusic(id)
	if err != nil {
		return err
	}

	if existingMusic.UserID != userID {
		return errors.New("non autorisé")
	}

	music.UpdatedAt = time.Now()
	return s.db.Model(existingMusic).Updates(music).Error
}

func (s *MusicService) DeleteMusic(id string, userID string) error {
	music, err := s.GetMusic(id)
	if err != nil {
		return err
	}

	if music.UserID != userID {
		return errors.New("Non autorisé")
	}
	return s.db.Delete(&models.Music{}, "id = ?", id).Error
}

func (s *MusicService) GetUserMusic(userID string) ([]models.Music, error) {
	var musics []models.Music
	if err := s.db.Where("user_id = ?", userID).Find(&musics).Error; err != nil {
		return nil, err
	}
	return musics, nil
}
