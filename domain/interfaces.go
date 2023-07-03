package domain

import (
	"context"
)

type AuthService interface {
	Authenticate(context.Context, AuthParams) (AuthTokenID, error)
	Validate(context.Context, AuthTokenID) error
	Expire(context.Context, AuthTokenID) error
}

type CredentialService interface {
	SaveCredential(context.Context, SaveParams) error
	ListCredentials(context.Context) ([]Credential, error)
}
