package repository

import "errors"

var (
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrReceiverNotFound        = errors.New("receiver not found")
	ErrTransactionFailed       = errors.New("transaction failed")
	ErrPurchaseFailed          = errors.New("purchase failed")
	ErrUserCreationFailed      = errors.New("user creation failed")
	ErrUserGettingFailed       = errors.New("user getting failed")
	ErrInfoGettingFailed       = errors.New("info getting failed")
	ErrNoMerchItemsFound       = errors.New("no merch items found")
	ErrMerchItemsGettingFailed = errors.New("merch items getting failed")
	ErrMerchItemScan           = errors.New("error scanning merch item")
)
