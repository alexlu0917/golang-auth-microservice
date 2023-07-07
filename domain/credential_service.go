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

func (s StCredentialService) ListCredentials(ctx context.Context) ([]Credential, error) {
	list, err := s.storage.FindAll(ctx)
	if err != nil {
		log.Printf("error: find all credentials %v\n", err)
		return list, fmt.Errorf("find all credentials failed")
	}

	return list, nil
}

func (s StCredentialService) RemoveCredential(ctx context.Context, id CredentialID) error {
	if err := s.storage.DeleteByID(ctx, id); err != nil {
		log.Printf("error: remove core %v\n", err)
		return fmt.Errorf("remove credential failed")
	}

	return nil
}
