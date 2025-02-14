package repository

import (
	"context"
	"database/sql"
	"errors"

	"merch-shop/internal/api"
	"merch-shop/pkg/database"
	"merch-shop/pkg/logger"
)

type InfoRepository struct {
	db database.DB
}

func NewInfoRepository(db database.DB) *InfoRepository {
	return &InfoRepository{
		db: db,
	}
}

func (i *InfoRepository) GetInfo(ctx context.Context, userID int) (*api.InfoResponse, error) {
	logger.Debug("InfoRepository.GetInfo: ", "message", "retrieving info", "userID", userID)

	query := `
	with
	user_balance as (
		select coins_balance 
		from users 
		where id = $1
	),
	user_inventory as (
		select mi.item_name, sum(p.quantity) as quantity
		from purchases p
		join merch_items mi on mi.id = p.item_id
		where p.user_id = $1
		group by mi.item_name
	),
	coin_history_received as (
		select u.username as to_user, t.amount as amount
		from transactions t
		join users u on u.id = t.receiver_id
		where t.sender_id = $1
	),
	coin_history_sent as (
		select u.username as from_user, t.amount as amount 
		from transactions t
		join users u on u.id = t.sender_id
		where t.receiver_id = $1
	),
	select
		(select coins_balance from user_balance) as coins_balance,
		coalesce(json_agg(json_build_object('type', i.item_type, 'quantity', i.quantity))filter (where i.item_type is not null), '[]'::json) as inventory,
		coalesce(json_agg(json_build_object('fromUser', r.from_user, 'amount', r.amount)) filter (where r.from_user is not null), '[]'::json) as coin_history_received,
		coalesce(json_agg(json_build_object('toUser', s.to_user, 'amount', s.amount)) filter (where s.to_user is not null), '[]'::json) as coin_history_sent
	from user_inventory as i
	full join coin_history_received as r on true
	full join coin_history_sent as r on true;
	`

	var info api.InfoResponse

	err := i.db.QueryRow(ctx, query, userID).Scan(&info.Coins, &info.Inventory, &info.CoinHistory.Sent, &info.CoinHistory.Received)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("InfoRepository.GetInfo: ", "message", "user not found", "userID", userID)
			return nil, ErrUserNotFound
		}
		logger.Error("InfoRepository.GetInfo: ", "message", "query execution error", "error", err, "userID", userID)
		return nil, ErrInfoGettingFailed
	}

	return &info, nil
}
