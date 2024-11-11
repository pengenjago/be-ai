package models

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID          string `gorm:"primaryKey"`
	Name        string
	Email       string
	Password    string
	NoTelephone string
	Role        string
}

func (m *User) BeforeCreate(db *gorm.DB) error {
	if m.ID == "" {
		m.ID = ulid.Make().String()
	}

	return nil
}
