package service

import (
	"encoding/base64"
	"errors"
	"log"

	"github.com/countenum404/Veksel/internal/repository"
	"github.com/countenum404/Veksel/internal/types"
)

type DefaultUserService struct {
	repo repository.UserRepository
}

func NewDefaultUserService(repos repository.UserRepository) *DefaultUserService {
	return &DefaultUserService{repo: repos}
}

func (dus *DefaultUserService) Authenticate(username, password string) (*types.User, error) {
	user, err := dus.repo.GetUser(username)
	if err != nil {
		return nil, errors.New("INVALID USERNAME OR PASSWORD")
	}
	userPassword, err := base64.StdEncoding.DecodeString(user.Password)
	if err != nil {
		log.Println(err)
	}
	if string(userPassword) != password {
		return nil, errors.New("INVALID USERNAME OR PASSWORD")
	}
	return user, nil
}
