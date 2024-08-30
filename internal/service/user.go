package service

import (
	"encoding/base64"
	"errors"

	"github.com/countenum404/Veksel/internal/repository"
	"github.com/countenum404/Veksel/internal/types"
	"github.com/countenum404/Veksel/pkg/logger"
)

type DefaultUserService struct {
	repo repository.UserRepository
}

func NewDefaultUserService(repos repository.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repos}
}

func (dus *DefaultUserService) GetUser(username, password string) (*types.User, error) {
	user, err := dus.repo.GetUser(username)
	if err != nil {
		return nil, errors.New("INVALID USERNAME OR PASSWORD")
	}
	userPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		logger.GetLogger().Err(err)
	}
	if string(userPassword) != password {
		return nil, errors.New("INVALID USERNAME OR PASSWORD")
	}
	return user, nil
}
