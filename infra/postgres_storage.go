package infra

import (
	"microauth/domain"

	"gorm.io/gorm"
)

type StPostgresStorage struct {
	*gorm.DB
}

type StCredential struct {
	gorm.Model
	Name      string `gorm:"unique"`
	Password  string
	AuthToken domain.AuthToken `gorm:"embedded;embeddedPrefix:auth_token_"`
}
