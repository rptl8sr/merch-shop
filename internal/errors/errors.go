package errors

import "errors"

var (
	ErrInvalidCredentials      = errors.New("invalid credentials")
	ErrHashingPassword         = errors.New("error hashing password")
	ErrNoMerchItemFound        = errors.New("no merch item found")
	ErrTransactionFailed       = errors.New("transaction failed")
	ErrPurchaseFailed          = errors.New("purchase failed")
	ErrInsufficientFunds       = errors.New("insufficient funds")
	ErrReceiverNotFound        = errors.New("receiver not found")
	ErrUserCreationFailed      = errors.New("user creation failed")
	ErrUserGettingFailed       = errors.New("user getting failed")
	ErrUserNotFound            = errors.New("user not found")
	ErrInfoGettingFailed       = errors.New("info getting failed")
	ErrNoMerchItemsFound       = errors.New("no merch items found")
	ErrMerchItemsGettingFailed = errors.New("merch items getting failed")
	ErrMerchItemScan           = errors.New("error scanning merch item")
	ErrReceiverSenderAreSame   = errors.New("sender and receiver are the same")
	ErrEmptyPassword           = errors.New("password cannot be empty")
)
