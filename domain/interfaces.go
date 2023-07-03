package domain

import (
	"context"
)

type AuthService interface {
	Authenticate(context.Context, AuthParams) (AuthTokenID, error)
}
