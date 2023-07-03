package domain

import (
	"context"
	"fmt"
	"log"
)

type StAuthService struct {
	storage Storage
	hasher  Hasher
}

func NewAuthService(s Storage, h Hasher) *StAuthService {
	return &StAuthService{storage: s, hasher: h}
}

func (s StAuthService) Authenticate(ctx context.Context, r AuthParams) (AuthTokenID, error) {
	a, err := s.storage.FindByName(ctx, r.Name)
	if err != nil {
		log.Printf("error: find credential by name %v\n", err)
		return "", fmt.Errorf("find credential by name failed")
	}

	if !s.hasher.Compare(r.Password, a.PasswordHash) {
		log.Println("error: Invalid Password")
		return "", fmt.Errorf("InValid Password")
	}

	return a.AuthToken.ID, nil
}
