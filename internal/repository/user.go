package repository

import (
	"context"
	"database/sql"
	"errors"

	internalErrors "merch-shop/internal/errors"
	"merch-shop/internal/model"
	"merch-shop/pkg/database"
	"merch-shop/pkg/logger"
)

type UserRepository struct {
	db           database.DB
	defaultCoins int
}

func NewUserRepository(db database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (u *UserRepository) GetOrCreateUser(ctx context.Context, username, passwordHash string) (*model.User, error) {
	logger.Debug("UserRepository.GetOrCreateUser: ", "message", "retrieving or creating user", "username", username)

	query := `
	with 
	user_upsert as (
		insert into users (username, password_hash, coins_balance) 
		values ($1, $2, $3) 
		on conflict (username) do nothing
		returning id, username, password_hash, coins_balance	
	),
	user_exists as (
		select id, username, password_hash, coins_balance 
		from users 
		where username = $1
	),
	transaction_insert as (
		insert into transactions (receiver_id, amount)
		select id, $3 from user_upsert
		returning id
	)
	select id, username, password_hash, coins_balance 
	from user_upsert
	union all
	select id, username, password_hash, coins_balance
	from user_exists
	limit 1;
	`

	var user model.User
	err := u.db.QueryRow(ctx, query, username, passwordHash, u.defaultCoins).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CoinsBalance)

	if err != nil {
		logger.Error("UserRepository.GetOrCreateUser: ", "message", "query execution error", "error", err, "username", username, "passwordHash", passwordHash)
		return nil, internalErrors.ErrUserCreationFailed
	}

	return &user, nil
}

func (u *UserRepository) GetUserByName(ctx context.Context, username string) (*model.User, error) {
	logger.Debug("UserRepository.GetUserByName: ", "message", "retrieving user by name", "username", username)

	var user model.User
	err := u.db.QueryRow(ctx, "select id, username, password_hash, coins_balance from users where username = $1", username).
		Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CoinsBalance)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Error("UserRepository.GetUserByName: ", "message", "user not found", "username", username)
			return nil, internalErrors.ErrUserNotFound
		}

		logger.Error("UserRepository.GetUserByName: ", "message", "query execution error", "error", err, "username", username)
		return nil, internalErrors.ErrUserGettingFailed
	}

	return &user, nil
}
