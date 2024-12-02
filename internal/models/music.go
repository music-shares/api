package models

import "time"

type Music struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"not null"`
	UserID    string    `json:"user_id" gorm:"not null"`
	AudioFile string    `json:"audio_file"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID"` // Relation avec User
}
