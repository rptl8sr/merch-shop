package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"merch-shop/internal/model"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type userService struct {
	repo UserRepository
}

type UserRepository interface {
	CreateUser(ctx context.Context, username, passwordHash string) (*model.User, error)
	GetUserByName(ctx context.Context, username string) (*model.User, error)
}

type UserService interface {
	GetOrCreate(ctx context.Context, username, password string) (*model.User, error)
	createUser(ctx context.Context, username, password string) (*model.User, error)
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) GetOrCreate(ctx context.Context, username, password string) (*model.User, error) {
	user, err := s.repo.GetUserByName(ctx, username)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			return s.createUser(ctx, username, password)
		}
	}

	if !comparePassword(user.PasswordHash, password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) createUser(ctx context.Context, username string, password string) (*model.User, error) {
	hash, err := hashPassword(password)
	if err != nil {
		return nil, err
	}

	user, err := s.repo.CreateUser(ctx, username, hash)
	if err != nil {
		return nil, err
	}

	// TODO: add default balance with tx

	return user, nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func comparePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
