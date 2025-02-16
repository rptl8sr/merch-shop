package service

import "context"

type TransactionService struct {
	repo TransactionRepository
}

type TransactionRepository interface {
	SendCoins(ctx context.Context, sender *int, receiver string, amount int) error
}

func NewTransactionService(repo TransactionRepository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (t *TransactionService) SendCoins(ctx context.Context, userID uint, toUser string, amount int) error {
	uID := int(userID)
	return t.repo.SendCoins(ctx, &uID, toUser, amount)
}
