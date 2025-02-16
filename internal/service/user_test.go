package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	internalErrors "merch-shop/internal/errors"
	"merch-shop/internal/model"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetOrCreateUser(ctx context.Context, username, passwordHash string) (*model.User, error) {
	args := m.Called(ctx, username, passwordHash)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(ctx, username)
	return args.Get(0).(*model.User), args.Error(1)
}

func TestUserService_GetOrCreate(t *testing.T) {
	ctx := context.Background()
	username := "test_user"
	password := "secure_password"

	// Генерируем хеш один раз для всех тестов
	passwordHash, err := hashPassword(password)
	assert.NoError(t, err)

	t.Run("successful user creation", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		expectedUser := &model.User{
			Username:     username,
			PasswordHash: passwordHash,
		}

		// Используем mock.Anything для passwordHash, так как соль генерируется случайно
		mockRepo.On("GetOrCreateUser", ctx, username, mock.AnythingOfType("string")).Return(expectedUser, nil)

		user, err := service.GetOrCreate(ctx, username, password)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		wrongPassword := "wrong_password"

		existingUser := &model.User{
			Username:     username,
			PasswordHash: passwordHash,
		}

		mockRepo.On("GetOrCreateUser", ctx, username, mock.Anything).Return(existingUser, nil)

		user, err := service.GetOrCreate(ctx, username, wrongPassword)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.ErrorIs(t, err, internalErrors.ErrInvalidCredentials)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		expectedError := errors.New("database error")
		mockRepo.On("GetOrCreateUser", ctx, username, mock.Anything).Return((*model.User)(nil), expectedError)

		user, err := service.GetOrCreate(ctx, username, password)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, expectedError, err)
	})

	t.Run("empty password handling", func(t *testing.T) {
		mockRepo := new(MockUserRepository)
		service := NewUserService(mockRepo)

		_, err := service.GetOrCreate(ctx, username, "")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
		mockRepo.AssertNotCalled(t, "GetOrCreateUser")
	})
}
