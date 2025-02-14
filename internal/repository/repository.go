package repository

import (
	"merch-shop/pkg/database"
)

type Repository struct {
	InfoRepository        *InfoRepository
	MerchRepository       *MerchRepository
	PurchaseRepository    *PurchaseRepository
	TransactionRepository *TransactionRepository
	UserRepository        *UserRepository
}

func NewRepository(db database.DB) *Repository {
	return &Repository{
		InfoRepository:        NewInfoRepository(db),
		MerchRepository:       NewMerchRepository(db),
		PurchaseRepository:    NewPurchaseRepository(db),
		TransactionRepository: NewTransactionRepository(db),
		UserRepository:        NewUserRepository(db),
	}
}
