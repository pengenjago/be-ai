package services

import (
	"be-ai/internal/constants"
	"be-ai/internal/dto"
	"be-ai/internal/models"
	"be-ai/internal/repositories"
	"be-ai/internal/token"
	"be-ai/util"
	"errors"
	"log"
	"time"
)

type UserService interface {
	CreateUser(param dto.UserReq) error
	LoginUser(param dto.UserLogin) (*dto.LoginRes, error)
}

var us *userServiceImpl

type userServiceImpl struct {
	userRepo repositories.UserRepository
}

func GetUserService() UserService {
	if us == nil {
		us = &userServiceImpl{
			userRepo: repositories.GetUserRepo(),
		}
	}
	return us
}

// ------------------------------------------

func (s *userServiceImpl) CreateUser(param dto.UserReq) error {

	hash, _ := util.HashPassword(param.Password)
	data := models.User{
		Name:        param.Name,
		Email:       param.Email,
		NoTelephone: param.NoTelephone,
		Password:    hash,
		Role:        constants.CUSTOMER,
	}

	exist := s.userRepo.GetUserByEmail(param.Email)
	if exist != nil {
		return errors.New("email already exists")
	}

	err := s.userRepo.CreateUser(&data)
	if err != nil {
		log.Println("failed to create user :", err.Error())
		return errors.New("failed to create user")
	}

	return nil
}

func (s *userServiceImpl) LoginUser(param dto.UserLogin) (*dto.LoginRes, error) {
	data := s.userRepo.GetUserByEmail(param.Email)
	if data == nil {
		return nil, constants.ErrInvalidUser
	}

	if err := util.CheckPassword(param.Password, data.Password); err != nil {
		return nil, constants.ErrInvalidUser
	}

	accessToken, accessPayload, err := token.NewJWT().Create(data.ID, data.Email, data.Role, 8*time.Hour)
	if err != nil {
		log.Println("error when create JWT :", err.Error())
		return nil, constants.ErrFailedLogin
	}

	refreshToken, refreshPayload, err := token.NewJWT().Create(data.ID, data.Email, data.Role, 16*time.Hour)
	if err != nil {
		log.Println("error when create JWT :", err.Error())
		return nil, constants.ErrFailedLogin
	}

	return &dto.LoginRes{
		AccessToken:                 accessToken,
		AccessTokenExpiredAt:        accessPayload.ExpiredAt,
		RefreshAccessToken:          refreshToken,
		RefreshAccessTokenExpiredAt: refreshPayload.ExpiredAt,
		UserRes: dto.UserRes{
			Name: data.Name,
			Role: data.Role,
		},
	}, nil
}
