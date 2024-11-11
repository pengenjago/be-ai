package repositories

import (
	"be-ai/config"
	"be-ai/internal/models"
	"gorm.io/gorm"
)

type AssistantsRepository interface {
	GetAll() []models.Assistants
	Create(data *models.Assistants) error
}

type assistantsRepoImpl struct {
	conn *gorm.DB
}

func GetAssistantsRepo() AssistantsRepository {
	return &assistantsRepoImpl{config.GetDb()}
}

// ------------------------------------------

func (a *assistantsRepoImpl) GetAll() []models.Assistants {
	var data []models.Assistants

	a.conn.Find(&data)

	return data
}

func (a *assistantsRepoImpl) Create(data *models.Assistants) error {
	return a.conn.Create(data).Error
}
