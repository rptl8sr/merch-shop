package service

import (
	"context"

	internalErrors "merch-shop/internal/errors"
	"merch-shop/internal/model"
	"merch-shop/pkg/logger"
)

type UserService struct {
	repo UserRepository
}

type UserRepository interface {
	GetOrCreateUser(ctx context.Context, username, passwordHash string) (*model.User, error)
	GetUserByName(ctx context.Context, username string) (*model.User, error)
}

func NewUserService(repo UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetOrCreate(ctx context.Context, username, password string) (*model.User, error) {
	passwordHash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.GetOrCreateUser(ctx, username, passwordHash)
	if err != nil {
		return nil, err
	}

	if !comparePassword(user.PasswordHash, password) {
		logger.Error("UserService.GetOrCreate: ", "message", "invalid credentials", "username", username)
		return nil, internalErrors.ErrInvalidCredentials
	}

	return user, nil
}
