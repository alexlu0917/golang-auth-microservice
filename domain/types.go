package domain

type AuthTokenID string

type AuthParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
