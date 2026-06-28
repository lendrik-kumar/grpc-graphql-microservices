package account

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	PostAccount(ctx context.Context, account *Account) (*Account, error)
	GetAccount(ctx context.Context, id string) (*Account, error)
	GetAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error)
}

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type accountService struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &accountService{repository: r}
}

func (s *accountService) PostAccount(ctx context.Context, account *Account) (*Account, error) {
	if account == nil {
		account = &Account{}
	}
	if account.ID == "" {
		account.ID = uuid.New().String()
	}

	if err := s.repository.PutAccount(ctx, account); err != nil {
		return nil, err
	}
	return account, nil
}

func (s *accountService) GetAccount(ctx context.Context, id string) (*Account, error) {
	return s.repository.GetAccountByID(ctx, id)
}

func (s *accountService) GetAccounts(ctx context.Context, skip uint64, take uint64) ([]*Account, error) {
	if take > 100 || (skip == 0 && take == 0) {
		take = 100
	}
	return s.repository.ListAccounts(ctx, skip, take)
}
