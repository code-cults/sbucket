package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type File struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
	BucketID  uuid.UUID `gorm:"type:uuid;not null"`
	FileName  string    `gorm:"not null"`
	Size      int64     `gorm:"not null"`
	Hash      string    `json:"hash"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	return nil
}