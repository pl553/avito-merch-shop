package service

import (
	"context"
	"fmt"
	"time"

	"merchshop/internal/model"
	"merchshop/internal/repository"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	jwtSecretKey  = "wqjeklasjdasnj"
	jwtExpiration = 12 * time.Hour
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

type MerchService struct {
	repo repository.MerchRepository
}

func NewMerchService(repo repository.MerchRepository) *MerchService {
	return &MerchService{repo: repo}
}

func (s *MerchService) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		if err == model.ErrUserNotFound {
			hashed, err := hashPassword(password)
			if err != nil {
				return "", fmt.Errorf("failed to hash password: %w", err)
			}
			if err := s.repo.CreateUser(ctx, username, hashed); err != nil {
				return "", fmt.Errorf("failed to create user: %w", err)
			}
			user, err = s.repo.GetUser(ctx, username)
			if err != nil {
				return "", fmt.Errorf("failed to retrieve created user: %w", err)
			}
		} else {
			return "", fmt.Errorf("failed to get user: %w", err)
		}
	}

	if !checkPasswordHash(password, user.PasswordHash) {
		return "", model.ErrInvalidPassword
	}

	claims := jwt.MapClaims{
		"sub": user.Username,
		"exp": time.Now().Add(jwtExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (s *MerchService) BuyItem(ctx context.Context, username, item string) error {
	price, err := s.repo.GetProductPrice(ctx, item)
	if err != nil {
		return fmt.Errorf("failed to get product price: %w", err)
	}

	err = s.repo.Atomic(ctx, func(r repository.MerchRepository) error {
		if err := r.DeductCoins(ctx, username, int32(price)); err != nil {
			return fmt.Errorf("failed to deduct coins: %w", err)
		}
		if cpRepo, ok := r.(interface {
			CreatePurchase(ctx context.Context, username string, item string, price int32) error
		}); ok {
			if err := cpRepo.CreatePurchase(ctx, username, item, int32(price)); err != nil {
				return fmt.Errorf("failed to create purchase record: %w", err)
			}
		} else {
			fmt.Println("warning: CreatePurchase not implemented in repository")
		}
		return nil
	})
	return err
}

func (s *MerchService) GetInfo(ctx context.Context, username string) (*model.Info, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	inv, err := s.repo.GetInventory(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}
	sent, err := s.repo.GetCoinHistorySent(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent coin history: %w", err)
	}
	received, err := s.repo.GetCoinHistoryReceived(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("failed to get received coin history: %w", err)
	}
	info := &model.Info{
		Coins:     user.Coins,
		Inventory: inv,
		CoinHistory: model.CoinHistory{
			Sent:     sent,
			Received: received,
		},
	}
	return info, nil
}

// SendCoin transfers coins from one user to another atomically.
func (s *MerchService) SendCoin(ctx context.Context, fromUsername, toUsername string, amount int32) error {
	return s.repo.Atomic(ctx, func(r repository.MerchRepository) error {
		if err := r.DeductCoins(ctx, fromUsername, amount); err != nil {
			return fmt.Errorf("failed to deduct coins from sender: %w", err)
		}
		if err := r.AddCoins(ctx, toUsername, amount); err != nil {
			return fmt.Errorf("failed to add coins to receiver: %w", err)
		}
		if err := r.InsertCoinTransfer(ctx, fromUsername, toUsername, amount); err != nil {
			return fmt.Errorf("failed to log coin transfer: %w", err)
		}
		return nil
	})
}
