package repository

import (
	"context"

	"merch-shop/pkg/database"
	"merch-shop/pkg/logger"
)

type TransactionRepository struct {
	db database.DB
}

func NewTransactionRepository(db database.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

func (t *TransactionRepository) SendCoins(ctx context.Context, sender *int, receiver string, amount int) error {
	logger.Debug("TransactionRepository.SendCoins: ", "message", "sending coins", "sender", *sender, "receiver", receiver, "amount", amount)

	query := `
	with
	receiver_check as (
		select id from users 
		where username = $2
	),
	sender_update as (
		update users
		set coins_balance = coins_balance - $3
		where id = $1 and coins_balance >= $3
		returning id	
	),
	receiver_update as (
		update users
		set coins_balance = coins_balance + $3
		where username = $2
		returning id
	),
	transaction_insert as (
		insert into transactions (sender_id, receiver_id, amount)
		select $1, (select id from receiver_update), $3 
		from sender_update
		returning id
	)
	select 
		(select count(*) from receiver_check) as receiver_exists,
		(select count(*) from sender_update) as sender_updated,
		(select count(*) from receiver_update) as receiver_updated,
		(select count(*) from transaction_insert) as transaction_inserted;
	`

	var (
		receiverExists      int
		senderUpdated       int
		receiverUpdated     int
		transactionInserted int
	)

	err := t.db.QueryRow(ctx, query, sender, receiver, amount).
		Scan(&receiverExists, &senderUpdated, &receiverUpdated, &transactionInserted)

	if err != nil {
		logger.Error("TransactionRepository.SendCoins: ", "message", "query execution error", "error", err, "sender", *sender, "receiver", receiver, "amount", amount)
		return ErrTransactionFailed
	}

	if receiverExists == 0 {
		logger.Error("TransactionRepository.SendCoins: ", "message", "receiver not found", "receiver", receiver)
		return ErrReceiverNotFound
	}

	if senderUpdated == 0 {
		logger.Error("TransactionRepository.SendCoins: ", "message", "insufficient funds", "sender", *sender, "amount", amount)
		return ErrInsufficientFunds
	}

	if receiverUpdated == 0 {
		logger.Error("TransactionRepository.SendCoins: ", "message", "receiver not updated", "receiver", receiver, "sender", *sender, "amount", amount)
		return ErrTransactionFailed
	}

	if transactionInserted == 0 {
		logger.Error("TransactionRepository.SendCoins: ", "message", "transaction not inserted", "receiver", receiver, "sender", *sender, "amount", amount)
		return ErrTransactionFailed
	}

	return nil
}
