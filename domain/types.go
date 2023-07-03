package domain

import "time"

type AuthTokenID string
type CredentialID uint
type AuthToken struct {
	ID        AuthTokenID
	ExpiresAt time.Time
}

type AuthParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SaveParams struct {
	ID       CredentialID `json:"id"`
	Name     string       `json:"name"`
	Password string       `json:"password"`
}

type Credential struct {
	ID           CredentialID `json:"id"`
	Name         string       `json:"name"`
	PasswordHash string       `json:"-"`
	AuthToken    AuthToken    `json:"-"`
}
