package model

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrItemNotFound      = errors.New("item not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type CoinTransferTo struct {
	ToUsername string
	Amount     uint32
}

type CoinTransferFrom struct {
	FromUsername string
	Amount       uint32
}

type InventoryItem struct {
	Item   string
	Amount uint32
}

type Purchase struct {
	Username string
	Item     string
	Price    uint32
}

type User struct {
	Username     string
	PasswordHash string
	Coins        uint32
}

type CoinHistory struct {
	Sent     []CoinTransferTo
	Received []CoinTransferFrom
}

type Info struct {
	Coins       uint32
	Inventory   []InventoryItem
	CoinHistory CoinHistory
}
