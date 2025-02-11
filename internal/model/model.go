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

type InfoResponse struct {
	Coins       int             `json:"coins"`
	Inventory   []InventoryItem `json:"inventory"`
	CoinHistory CoinHistory     `json:"coinHistory"`
}

type InventoryItem struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type CoinHistory struct {
	Received []ReceivedTransaction `json:"received"`
	Sent     []SentTransaction     `json:"sent"`
}

type ReceivedTransaction struct {
	FromUser User `json:"fromUser"`
	Amount   int  `json:"amount"`
}

type SentTransaction struct {
	ToUser User `json:"toUser"`
	Amount int  `json:"amount"`
}

type ErrorResponse struct {
	Errors string `json:"errors"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type SendCountRequest struct {
	ReceiverID int `json:"toUser"`
	Amount     int `json:"amount"`
}
