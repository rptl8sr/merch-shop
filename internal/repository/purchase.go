package repository

import (
	"context"
	"merch-shop/pkg/logger"

	"merch-shop/internal/model"
	"merch-shop/pkg/database"
)

type PurchaseRepository struct {
	db database.DB
}

func NewPurchaseRepository(db database.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (p *PurchaseRepository) MakePurchase(ctx context.Context, user model.User, merch model.MerchItem, quantity int) error {
	logger.Debug("PurchaseRepository.MakePurchase: ", "message", "making purchase", "user", user, "merch", merch, "quantity", quantity)

	query := `
	with
	user_check as (
		select id 
		from users
		where id = $1
		for update
	),
	user_coins_balance as (
		update users 
		set coins_balance = coins_balance - $3 * $4
		where id = $1 and coins_balance >= $3 * $4
		returning id
	),
	purchase_insert as (
		insert into purchases (user_id, item_id, quantity, total_cost)
		select id, $2, $4, $3 * $4 
		from user_coins_balance
		returning id
	)
	select 
		(select count(*) from user_check) as user_exists,
		(select count(*) from user_coins_balance) as balance_updated,
		(select count(*) from purchase_insert) as purchase_inserted;
	`

	var (
		userExists       int
		balanceUpdated   int
		purchaseInserted int
	)

	err := p.db.QueryRow(ctx, query, user.ID, merch.ID, merch.Price, quantity).
		Scan(&userExists, &balanceUpdated, &purchaseInserted)

	if err != nil {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "query execution error", "error", err, "user", user, "merch", merch, "quantity", quantity)
		return ErrPurchaseFailed
	}

	if userExists == 0 {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "user not found", "user", user)
		return ErrUserNotFound
	}

	if balanceUpdated == 0 {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "insufficient funds", "user", user, "merch", merch, "quantity", quantity)
		return ErrInsufficientFunds
	}

	if purchaseInserted == 0 {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "purchase not inserted", "user", user, "merch", merch, "quantity", quantity)
		return ErrPurchaseFailed
	}

	return nil
}
