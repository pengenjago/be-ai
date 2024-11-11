package repositories

import (
	"be-ai/config"
	"be-ai/internal/models"
	"gorm.io/gorm"
)

type ThreadRepository interface {
	GetAll(userId string) []models.Thread
	Create(userId, threadId, assistantId string) error
}

type threadRepoImpl struct {
	conn *gorm.DB
}

func GetThreadRepo() ThreadRepository {
	return &threadRepoImpl{config.GetDb()}
}

// ------------------------------------------

func (t *threadRepoImpl) Create(userId, threadId, assistantId string) error {
	data := models.Thread{
		ID:          threadId,
		UserId:      userId,
		AssistantId: assistantId,
	}
	return t.conn.Create(&data).Error
}

func (t *threadRepoImpl) GetAll(userId string) []models.Thread {
	var data []models.Thread

	t.conn.Where("user_id = ?", userId).Debug().Find(&data)

	return data
}
