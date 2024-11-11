package models

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Assistants struct {
	gorm.Model
	ID           string `gorm:"primaryKey"`
	Name         string
	Instructions string
	GptModel     string
	VectorID     string
}

func (m *Assistants) BeforeCreate(db *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}

	return nil
}
