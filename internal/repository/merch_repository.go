package repository

import (
	"context"
	"merchshop/internal/model"
)

type MerchRepository interface {
	Atomic(context.Context, func(r MerchRepository) error) error
	CreateUser(ctx context.Context, username string, passwordHash string) error
	AddCoins(ctx context.Context, username string, amount int32) error
	DeductCoins(ctx context.Context, username string, amount int32) error
	InsertCoinTransfer(ctx context.Context, fromUsername string, toUsername string, amount int32) error
	GetCoinHistorySent(ctx context.Context, username string) ([]model.CoinTransferTo, error)
	GetCoinHistoryReceived(ctx context.Context, username string) ([]model.CoinTransferFrom, error)
	GetInventory(ctx context.Context, username string) ([]model.InventoryItem, error)
	GetProductPrice(ctx context.Context, item string) (uint32, error)
	GetUser(ctx context.Context, username string) (*model.User, error)
}
