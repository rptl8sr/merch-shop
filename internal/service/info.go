package service

import (
	"context"

	"merch-shop/internal/api"
)

type InfoService struct {
	repo InfoRepository
}

type InfoRepository interface {
	GetInfo(ctx context.Context, userID int) (*api.InfoResponse, error)
}

func NewInfoService(repo InfoRepository) *InfoService {
	return &InfoService{repo: repo}
}

func (i *InfoService) GetUserInfo(ctx context.Context, userID uint) (*api.InfoResponse, error) {
	uID := int(userID)
	return i.repo.GetInfo(ctx, uID)
}
