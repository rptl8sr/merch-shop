package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) SendCoins(ctx context.Context, sender *int, receiver string, amount int) error {
	args := m.Called(ctx, sender, receiver, amount)
	return args.Error(0)
}

func TestTransactionService_SendCoins(t *testing.T) {
	ctx := context.Background()
	userID := uint(123)
	receiver := "user456"
	amount := 100

	t.Run("successful coins transfer", func(t *testing.T) {
		mockRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockRepo)

		expectedSender := int(userID)
		mockRepo.On(
			"SendCoins",
			ctx,
			&expectedSender,
			receiver,
			amount,
		).Return(nil)

		err := service.SendCoins(ctx, userID, receiver, amount)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository returns error", func(t *testing.T) {
		mockRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockRepo)

		expectedError := errors.New("database error")
		mockRepo.On(
			"SendCoins",
			ctx,
			mock.AnythingOfType("*int"),
			receiver,
			amount,
		).Return(expectedError)

		err := service.SendCoins(ctx, userID, receiver, amount)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("correct parameters passed to repository", func(t *testing.T) {
		mockRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockRepo)

		var capturedSender *int
		mockRepo.On(
			"SendCoins",
			ctx,
			mock.MatchedBy(func(s *int) bool {
				capturedSender = s
				return true
			}),
			receiver,
			amount,
		).Return(nil)

		err := service.SendCoins(ctx, userID, receiver, amount)

		assert.NoError(t, err)
		assert.NotNil(t, capturedSender)
		assert.Equal(t, int(userID), *capturedSender)
	})

	t.Run("zero amount handling", func(t *testing.T) {
		mockRepo := new(MockTransactionRepository)
		service := NewTransactionService(mockRepo)

		mockRepo.On(
			"SendCoins",
			ctx,
			mock.AnythingOfType("*int"),
			receiver,
			0,
		).Return(nil)

		err := service.SendCoins(ctx, userID, receiver, 0)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "SendCoins", ctx, mock.Anything, receiver, 0)
	})
}
