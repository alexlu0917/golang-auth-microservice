package infra

import "golang.org/x/crypto/bcrypt"

type BCryptHasher struct{}

func (b BCryptHasher) Hash(pass string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b BCryptHasher) compare(pass string, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)) == nil
}
