package models

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type Thread struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	UserId      string
	AssistantId string
}

func (m *Thread) BeforeCreate(db *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}

	return nil
}
