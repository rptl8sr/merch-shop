package repository

import (
	"context"
	"merch-shop/internal/errors"

	"merch-shop/pkg/database"
	"merch-shop/pkg/logger"
)

type PurchaseRepository struct {
	db database.DB
}

func NewPurchaseRepository(db database.DB) *PurchaseRepository {
	return &PurchaseRepository{
		db: db,
	}
}

func (p *PurchaseRepository) MakePurchase(ctx context.Context, userID, merchID, price, quantity int) error {
	logger.Debug("PurchaseRepository.MakePurchase: ", "message", "making purchase", "userID", userID, "merchID", merchID, "quantity", quantity)

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

	err := p.db.QueryRow(ctx, query, userID, merchID, price, quantity).
		Scan(&userExists, &balanceUpdated, &purchaseInserted)

	if err != nil {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "query execution error", "error", err, "userID", userID, "merchID", merchID, "quantity", quantity)
		return errors.ErrPurchaseFailed
	}

	if userExists == 0 {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "user not found", "userID", userID)
		return errors.ErrUserNotFound
	}

	if balanceUpdated == 0 {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "insufficient funds", "userID", userID, "merchID", merchID, "quantity", quantity)
		return errors.ErrInsufficientFunds
	}

	if purchaseInserted == 0 {
		logger.Error("PurchaseRepository.MakePurchase: ", "message", "purchase not inserted", "userID", userID, "merchID", merchID, "quantity", quantity)
		return errors.ErrPurchaseFailed
	}

	return nil
}
