package repository

import (
	"context"
	"errors"

	"merchshop/internal/db/queries"
	"merchshop/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	uniqueViolationErrCode     = "23505"
	foreignKeyViolationErrCode = "23503"
)

type PgMerchRepository struct {
	queries *queries.Queries
	pool    *pgxpool.Pool
}

func NewPgMerchRepository(ctx context.Context, connString string) (*PgMerchRepository, error) {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}
	qs := queries.New(pool)
	return &PgMerchRepository{
		queries: qs,
		pool:    pool,
	}, nil
}

func (r *PgMerchRepository) Close() {
	r.pool.Close()
}

func (r *PgMerchRepository) Atomic(ctx context.Context, fn func(r MerchRepository) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	qtx := r.queries.WithTx(tx)
	txRepo := &PgMerchRepository{
		queries: qtx,
		pool:    r.pool,
	}
	if err := fn(txRepo); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return tx.Commit(ctx)
}

func (r *PgMerchRepository) CreateUser(ctx context.Context, username string, passwordHash string) error {
	err := r.queries.CreateUser(ctx, queries.CreateUserParams{
		Username:     username,
		PasswordHash: passwordHash,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationErrCode {
			return model.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (r *PgMerchRepository) AddCoins(ctx context.Context, username string, amount int32) error {
	rows, err := r.queries.AddCoins(ctx, queries.AddCoinsParams{
		Coins:    amount,
		Username: username,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *PgMerchRepository) DeductCoins(ctx context.Context, username string, amount int32) error {
	rows, err := r.queries.DeductCoins(ctx, queries.DeductCoinsParams{
		Coins:    amount,
		Username: username,
	})
	if err != nil {
		return err
	}
	if rows == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *PgMerchRepository) InsertCoinTransfer(ctx context.Context, fromUsername string, toUsername string, amount int32) error {
	err := r.queries.InsertCoinTransfer(ctx, queries.InsertCoinTransferParams{
		FromUsername: fromUsername,
		ToUsername:   toUsername,
		Amount:       amount,
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == foreignKeyViolationErrCode {
			return model.ErrUserNotFound
		}
		return err
	}
	return nil
}

func (r *PgMerchRepository) GetCoinHistorySent(ctx context.Context, username string) ([]model.CoinTransferTo, error) {
	rows, err := r.queries.GetCoinHistorySent(ctx, username)
	if err != nil {
		return nil, err
	}
	var transfers []model.CoinTransferTo
	for _, row := range rows {
		transfers = append(transfers, model.CoinTransferTo{
			ToUsername: row.ToUsername,
			Amount:     uint32(row.Amount),
		})
	}
	return transfers, nil
}

func (r *PgMerchRepository) GetCoinHistoryReceived(ctx context.Context, username string) ([]model.CoinTransferFrom, error) {
	rows, err := r.queries.GetCoinHistoryReceived(ctx, username)
	if err != nil {
		return nil, err
	}
	var transfers []model.CoinTransferFrom
	for _, row := range rows {
		transfers = append(transfers, model.CoinTransferFrom{
			FromUsername: row.FromUsername,
			Amount:       uint32(row.Amount),
		})
	}
	return transfers, nil
}

func (r *PgMerchRepository) GetInventory(ctx context.Context, username string) ([]model.InventoryItem, error) {
	rows, err := r.queries.ListInventory(ctx, username)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	var inventory []model.InventoryItem
	for _, row := range rows {
		inventory = append(inventory, model.InventoryItem{
			Item:   row.Item,
			Amount: uint32(row.Quantity),
		})
	}
	return inventory, nil
}

func (r *PgMerchRepository) GetProductPrice(ctx context.Context, item string) (uint32, error) {
	price, err := r.queries.GetProductPrice(ctx, item)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, model.ErrItemNotFound
		}
		return 0, err
	}
	return uint32(price), nil
}

func (r *PgMerchRepository) GetUser(ctx context.Context, username string) (*model.User, error) {
	user, err := r.queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}
	return &model.User{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Coins:        uint32(user.Coins),
	}, nil
}
