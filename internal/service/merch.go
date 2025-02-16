package service

import (
	"context"
	"sync"

	"merch-shop/internal/model"
	"merch-shop/internal/repository"
	"merch-shop/pkg/logger"
)

var (
	merchCache *MerchService
	merchOnce  sync.Once
)

type MerchServicer interface {
	GetMerchList(ctx context.Context) error
	GetMerchPrice(key string) int
	GetMerchItem(key string) (model.MerchItem, bool)
}

type MerchService struct {
	repo  MerchRepository
	items map[string]model.MerchItem
	sync.RWMutex
}

type MerchRepository interface {
	GetMerchList(ctx context.Context) ([]repository.MerchDTO, error)
}

func NewMerchService(repo MerchRepository) *MerchService {
	merchOnce.Do(func() {
		merchCache = &MerchService{
			repo:  repo,
			items: make(map[string]model.MerchItem),
		}
	})

	return merchCache
}

func (m *MerchService) GetMerchList(ctx context.Context) error {
	merchDTOs, err := m.repo.GetMerchList(ctx)
	if err != nil {
		return err
	}

	m.Lock()
	defer m.Unlock()

	for _, dto := range merchDTOs {
		if dto.Price > 0 && dto.ItemName != "" {
			m.items[dto.ItemName] = model.MerchItem{
				Meta:     model.Meta{ID: dto.ID},
				ItemName: dto.ItemName,
				Price:    dto.Price,
			}
		}
	}

	if len(m.items) == 0 {
		logger.Warn("MerchService.GetMerchList: no merch items found")
	}

	return nil
}

func (m *MerchService) GetMerchPrice(key string) int {
	m.RLock()
	defer m.RUnlock()

	return m.items[key].Price
}

func (m *MerchService) GetMerchItem(key string) (model.MerchItem, bool) {
	m.RLock()
	defer m.RUnlock()

	item, ok := m.items[key]
	return item, ok
}
