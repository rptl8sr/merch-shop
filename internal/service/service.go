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
	return &Services{
		InfoService:        NewInfoService(repo.InfoRepository),
		MerchService:       NewMerchService(repo.MerchRepository),
		PurchaseService:    NewPurchaseService(repo.PurchaseRepository),
		TransactionService: NewTransactionService(repo.TransactionRepository),
		UserService:        NewUserService(repo.UserRepository),
	}
}
