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
	credential, err := s.storage.FindByName(ctx, r.Name)
	if err != nil {
		log.Printf("error: find credential by name %v\n", err)
		return "", fmt.Errorf("find credential by name failed")
	}

	if !s.hasher.Compare(r.Password, credential.PasswordHash) {
		log.Println("error: Invalid Password")
		return "", fmt.Errorf("InValid Password")
	}

	if err = credential.GenerateAuthToken(); err != nil {
		log.Printf("error: save credential %v\n", err)
		return "", fmt.Errorf("generate auth token failed")
	}

	if err = s.storage.Save(ctx, credential); err != nil {
		log.Printf("error: save credential %v\n", err)
		return "", fmt.Errorf("save credential failed")
	}

	return credential.AuthToken.ID, nil
}

func (s StAuthService) Validate(ctx context.Context, id AuthTokenID) error {
	if id == "-" {
		return fmt.Errorf("token invalid")
	}

	credential, err := s.storage.FindByAuthTokenID(ctx, id)
	if err != nil {
		log.Printf("error: find credential by token id %v\n", err)
		return fmt.Errorf("find credential by token id failed")
	}

	if credential.AuthTokenExpired() {
		return fmt.Errorf("token expired")
	}

	return nil
}

func (s StAuthService) Expire(ctx context.Context, id AuthTokenID) error {
	credential, err := s.storage.FindByAuthTokenID(ctx, id)
	if err != nil {
		log.Printf("error: find credential by token id %v\n", err)
		return fmt.Errorf("find credential by token id failed")
	}

	credential.ExpireAuthToken()
	if err = s.storage.Save(ctx, credential); err != nil {
		log.Printf("error: saving credential on token expire %v\n", err)
		return fmt.Errorf("saving credential on token expire failed")
	}

	return nil
}
