package service

import (
	"merch-shop/internal/repository"
)

type Services struct {
	InfoService        *InfoService
	MerchService       *MerchService
	PurchaseService    *PurchaseService
	TransactionService *TransactionService
	UserService        *UserService
}

func NewServices(repo repository.Repository) *Services {
	merchService := NewMerchService(repo.MerchRepository)

	return &Services{
		InfoService:        NewInfoService(repo.InfoRepository),
		MerchService:       merchService,
		PurchaseService:    NewPurchaseService(repo.PurchaseRepository, merchService),
		TransactionService: NewTransactionService(repo.TransactionRepository),
		UserService:        NewUserService(repo.UserRepository),
	}
}
