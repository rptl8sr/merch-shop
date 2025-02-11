package model

import "time"

type Meta struct {
	ID        uint      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type User struct {
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CoinsBalance int    `json:"coins_balance"`
}

type MerchItem struct {
	Meta
	ItemName string `json:"name"`
}

type Transaction struct {
	Meta
	SenderID   *int `json:"sender_id"`
	ReceiverID int  `json:"receiver_id"`
	Amount     int  `json:"amount"`
}

type Purchase struct {
	Meta
	UserID    int `json:"user_id"`
	ItemID    int `json:"item_id"`
	Quantity  int `json:"quantity"`
	TotalCost int `json:"total_cost"`
}
