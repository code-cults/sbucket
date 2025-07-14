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
	Hash      string
	MimeType  string
	CreatedAt time.Time
}

func (f *File) BeforeCreate(tx *gorm.DB) (err error) {
	f.ID = uuid.New()
	return nil
}