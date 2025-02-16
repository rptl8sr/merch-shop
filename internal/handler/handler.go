package handler

import (
	"context"

	"merch-shop/internal/api"
	"merch-shop/internal/model"
	"merch-shop/internal/service"
)

var (
	defaultQuantity = 1
)

var (
	success = map[string]string{"status": "success"}
)

type InfoService interface {
	GetUserInfo(context context.Context, userID uint) (*api.InfoResponse, error)
}

type PurchaseService interface {
	BuyItem(ctx context.Context, userID uint, merchName string, quantity *int) error
}

type TransactionService interface {
	SendCoins(context context.Context, userID uint, toUser string, amount int) error
}

type UserService interface {
	GetOrCreate(ctx context.Context, username, password string) (*model.User, error)
}

type Handler struct {
	secret             string
	infoService        InfoService
	purchaseService    PurchaseService
	transactionService TransactionService
	userService        UserService
}

func New(secret string, services *service.Services) *Handler {
	return &Handler{
		secret:             secret,
		infoService:        services.InfoService,
		purchaseService:    services.PurchaseService,
		transactionService: services.TransactionService,
		userService:        services.UserService,
	}
}
