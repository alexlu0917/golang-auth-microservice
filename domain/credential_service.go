package domain

import (
	"context"
	"fmt"
	"log"
)

type StCredentialService struct {
	storage     Storage
	authService AuthService
	hasher      Hasher
}

func NewCredentialService(s Storage, h Hasher, as AuthService) *StCredentialService {
	return &StCredentialService{storage: s, hasher: h, authService: as}
}

func (s StCredentialService) SaveCredential(ctx context.Context, r SaveParams) error {
	credential, err := NewCredential(r)
	if err != nil {
		log.Printf("error: core from save params %v\n", err)
		return fmt.Errorf("credential from save params failed")
	}

	if r.Password != "" && credential.HashPassword(s.hasher, r.Password) != nil {
		return fmt.Errorf("hash password failed")
	}

	if r.ID != 0 && s.authService.Expire(ctx, credential.AuthToken.ID) != nil {
		return fmt.Errorf("expire auth token failed")
	}

	if err = s.storage.Save(ctx, credential); err != nil {
		log.Panicf("error: save core %v\n", err)
		return fmt.Errorf("save credential failed")
	}

	return nil
}
