package repositories

import (
	"be-ai/config"
	"be-ai/internal/dto"
	"be-ai/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	DeleteUser(user *models.User) error
	GetUserByEmail(email string) *models.User
	FindAll(query dto.UserQuery) ([]models.User, int)
	GetById(id string) *models.User
}

type userRepoImpl struct {
	conn *gorm.DB
}

func GetUserRepo() UserRepository {
	return &userRepoImpl{config.GetDb()}
}

// ------------------------------------------

func (u *userRepoImpl) CreateUser(user *models.User) error {
	return u.conn.Create(user).Error
}

func (u *userRepoImpl) UpdateUser(user *models.User) error {
	return u.conn.Save(user).Error
}

func (u *userRepoImpl) DeleteUser(user *models.User) error {
	return u.conn.Delete(user).Error
}

func (u *userRepoImpl) GetUserByEmail(email string) *models.User {
	var user models.User

	err := u.conn.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil
	}

	return &user
}
func (u *userRepoImpl) GetById(id string) *models.User {
	var user models.User

	err := u.conn.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil
	}

	return &user
}

func (u *userRepoImpl) FindAll(query dto.UserQuery) ([]models.User, int) {
	var user []models.User
	tx := u.conn

	if query.Search != "" {
		tx = tx.Where("name like ?", "%"+query.Search+"%")
	}

	total := tx.Find(&user).RowsAffected
	if query.PageSize > 0 {
		tx = tx.Offset((query.PageNo - 1) * query.PageSize).Limit(query.PageSize)
	}

	tx.Find(&user)

	return user, int(total)
}
