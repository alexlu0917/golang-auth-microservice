package infra

import (
	"context"
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

func toCredential(a *StCredential) domain.Credential {
	return domain.Credential{
		ID:           domain.CredentialID(a.ID),
		Name:         a.Name,
		PasswordHash: a.Password,
		AuthToken:    a.AuthToken,
	}
}

func fromCredential(a domain.Credential) *StCredential {
	return &StCredential{
		Name:      a.Name,
		Password:  a.PasswordHash,
		AuthToken: a.AuthToken,
	}
}

func NewPostgressStorage(db *gorm.DB) (*StPostgresStorage, error) {
	if err := db.AutoMigrate(&StCredential{}); err != nil {
		return &StPostgresStorage{}, err
	}

	return &StPostgresStorage{db}, nil
}

func (s StPostgresStorage) Save(ctx context.Context, a domain.Credential) error {
	row := fromCredential(a)
	if a.ID == domain.CredentialID(0) {
		return s.Create(row).Error
	}

	return s.WithContext(ctx).Model(&StCredential{}).Where("id = ?", a.ID).Updates(row).Error
}

func (s StPostgresStorage) FindAll(ctx context.Context) ([]domain.Credential, error) {
	var rows []StCredential
	var credentials []domain.Credential

	tx := s.WithContext(ctx).Find(&rows)
	if tx.Error != nil {
		return credentials, tx.Error
	}

	for _, row := range rows {
		credentials = append(credentials, toCredential(&row))
	}

	return credentials, nil
}

func (s StPostgresStorage) FindByID(ctx context.Context, id domain.CredentialID) (domain.Credential, error) {
	var row StCredential
	tx := s.First(ctx, &row, id)
	if tx.Error != nil {
		return toCredential(&row), tx.Error
	}

	return toCredential(&row), nil
}

func (s StPostgresStorage) FindByName(ctx context.Context, name string) (domain.Credential, error) {
	var row StCredential
	tx := s.WithContext(ctx).First(&row, "name = ?", name)
	if tx.Error != nil {
		return toCredential(&row), tx.Error
	}

	return toCredential(&row), nil
}

func (s StPostgresStorage) FindByAuthTokenID(ctx context.Context, id domain.AuthTokenID) (domain.Credential, error) {
	var row StCredential
	tx := s.WithContext(ctx).First(&row, "auth_token_id = ?", id)
	if tx.Error != nil {
		return toCredential(&row), tx.Error
	}

	return toCredential(&row), nil
}

func (s StPostgresStorage) DeleteByID(ctx context.Context, id domain.CredentialID) error {
	return s.WithContext(ctx).Delete(&StCredential{}, id).Error
}
