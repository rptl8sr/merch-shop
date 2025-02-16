package service

import (
	"context"

	internalErrors "merch-shop/internal/errors"
)

const (
	defaultQuantity = 1
)

type PurchaseService struct {
	repo         PurchaseRepository
	merchService MerchServicer
}

type PurchaseRepository interface {
	MakePurchase(ctx context.Context, userID, merchID, price, quantity int) error
}

func NewPurchaseService(repo PurchaseRepository, merchService MerchServicer) *PurchaseService {
	return &PurchaseService{
		repo:         repo,
		merchService: merchService,
	}
}

func (s *PurchaseService) BuyItem(ctx context.Context, userID uint, merchName string, quantity *int) error {
	merchItem, ok := s.merchService.GetMerchItem(merchName)
	if !ok {
		return internalErrors.ErrNoMerchItemFound
	}

	quantityToUse := defaultQuantity
	if quantity != nil {
		quantityToUse = *quantity
	}

	return s.repo.MakePurchase(ctx, int(userID), int(merchItem.ID), merchItem.Price, quantityToUse)
}
