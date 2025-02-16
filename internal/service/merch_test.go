package service

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"merch-shop/internal/model"
	"merch-shop/internal/repository"
	"merch-shop/pkg/logger"
)

type MockMerchRepository struct {
	mock.Mock
}

func (m *MockMerchRepository) GetMerchList(ctx context.Context) ([]repository.MerchDTO, error) {
	args := m.Called(ctx)
	return args.Get(0).([]repository.MerchDTO), args.Error(1)
}

func resetMerchSingleton() {
	merchCache = nil
	merchOnce = sync.Once{}
}

func TestMain(m *testing.M) {
	logger.Init(slog.LevelInfo)
	m.Run()
}

func TestMerchService_GetMerchList(t *testing.T) {
	ctx := context.Background()

	t.Run("successful data loading", func(t *testing.T) {
		resetMerchSingleton()
		mockRepo := new(MockMerchRepository)
		service := NewMerchService(mockRepo)

		mockData := []repository.MerchDTO{
			{ID: 1, ItemName: "T-Shirt", Price: 100},
			{ID: 2, ItemName: "Hat", Price: 50},
		}
		mockRepo.On("GetMerchList", ctx).Return(mockData, nil)

		err := service.GetMerchList(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 2, len(service.items))
		assert.Equal(t, 100, service.GetMerchPrice("T-Shirt"))
		assert.Equal(t, 50, service.GetMerchPrice("Hat"))
	})

	t.Run("empty response from repository", func(t *testing.T) {
		resetMerchSingleton()
		mockRepo := new(MockMerchRepository)
		service := NewMerchService(mockRepo)

		mockRepo.On("GetMerchList", ctx).Return([]repository.MerchDTO{}, nil)

		err := service.GetMerchList(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 0, len(service.items))
	})

	t.Run("repository returns error", func(t *testing.T) {
		resetMerchSingleton()
		mockRepo := new(MockMerchRepository)
		service := NewMerchService(mockRepo)

		mockRepo.On("GetMerchList", ctx).Return([]repository.MerchDTO{}, errors.New("database error"))

		err := service.GetMerchList(ctx)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		assert.Equal(t, 0, len(service.items))
	})

	t.Run("invalid item data handling", func(t *testing.T) {
		resetMerchSingleton()
		mockRepo := new(MockMerchRepository)
		service := NewMerchService(mockRepo)

		invalidData := []repository.MerchDTO{
			{ID: 1, ItemName: "", Price: 100},
			{ID: 2, ItemName: "Hat", Price: -50},
			{ID: 3, ItemName: "Valid", Price: 75},
		}
		mockRepo.On("GetMerchList", ctx).Return(invalidData, nil)

		err := service.GetMerchList(ctx)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(service.items))
		assert.Equal(t, 75, service.GetMerchPrice("Valid"))
		assert.Equal(t, 0, service.GetMerchPrice("Hat"))
	})
}

func TestMerchService_GetMerchPrice(t *testing.T) {
	resetMerchSingleton()
	mockRepo := new(MockMerchRepository)
	service := NewMerchService(mockRepo)

	service.items = map[string]model.MerchItem{
		"Jacket": {Price: 200},
	}

	t.Run("existing item price", func(t *testing.T) {
		assert.Equal(t, 200, service.GetMerchPrice("Jacket"))
	})

	t.Run("non-existent item price", func(t *testing.T) {
		assert.Equal(t, 0, service.GetMerchPrice("Gloves"))
	})
}

func TestMerchService_GetMerchItem(t *testing.T) {
	resetMerchSingleton()
	mockRepo := new(MockMerchRepository)
	service := NewMerchService(mockRepo)

	service.items = map[string]model.MerchItem{
		"Cap": {ItemName: "Cap", Price: 30},
	}

	t.Run("existing item retrieval", func(t *testing.T) {
		item, exists := service.GetMerchItem("Cap")
		assert.True(t, exists)
		assert.Equal(t, "Cap", item.ItemName)
		assert.Equal(t, 30, item.Price)
	})

	t.Run("non-existent item retrieval", func(t *testing.T) {
		_, exists := service.GetMerchItem("Scarf")
		assert.False(t, exists)
	})
}

func TestNewMerchService_Singleton(t *testing.T) {
	resetMerchSingleton()
	mockRepo1 := new(MockMerchRepository)
	mockRepo2 := new(MockMerchRepository)

	service1 := NewMerchService(mockRepo1)
	service2 := NewMerchService(mockRepo2)

	assert.Same(t, service1, service2, "Should return same instance for singleton")
	assert.Equal(t, mockRepo1, service1.repo, "Should keep first repository instance")
}
