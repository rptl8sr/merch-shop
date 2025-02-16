package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"merch-shop/internal/errors"
	"merch-shop/internal/model"
)

type MockPurchaseRepository struct {
	mock.Mock
}

func (m *MockPurchaseRepository) MakePurchase(ctx context.Context, userID, merchID, price, quantity int) error {
	args := m.Called(ctx, userID, merchID, price, quantity)
	return args.Error(0)
}

type MockMerchService struct {
	items map[string]model.MerchItem
}

func (m *MockMerchService) GetMerchItem(key string) (model.MerchItem, bool) {
	item, ok := m.items[key]
	return item, ok
}

func (m *MockMerchService) GetMerchPrice(key string) int {
	return m.items[key].Price
}

func (m *MockMerchService) GetMerchList(_ context.Context) error {
	return nil
}

func TestPurchaseService_BuyItem(t *testing.T) {
	ctx := context.Background()
	merchItem := model.MerchItem{
		Meta:     model.Meta{ID: 1},
		ItemName: "T-Shirt",
		Price:    100,
	}

	t.Run("successful purchase with custom quantity", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepository)
		mockMerch := &MockMerchService{
			items: map[string]model.MerchItem{
				"T-Shirt": merchItem,
			},
		}

		service := NewPurchaseService(mockRepo, mockMerch)

		qty := 3
		mockRepo.On("MakePurchase", ctx, 123, 1, 100, 3).Return(nil)

		err := service.BuyItem(ctx, 123, "T-Shirt", &qty)

		assert.NoError(t, err)
		mockRepo.AssertCalled(t, "MakePurchase", ctx, 123, 1, 100, 3)
	})

	t.Run("item not found", func(t *testing.T) {
		mockRepo := new(MockPurchaseRepository)
		mockMerch := &MockMerchService{
			items: make(map[string]model.MerchItem),
		}

		service := NewPurchaseService(mockRepo, mockMerch)

		err := service.BuyItem(ctx, 123, "Non-Existent", nil)

		assert.ErrorIs(t, err, errors.ErrNoMerchItemFound)
		mockRepo.AssertNotCalled(t, "MakePurchase")
	})
}
