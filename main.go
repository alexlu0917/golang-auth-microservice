package main

import (
	"context"
	"fmt"
	"log"
	"microauth/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

}

func mustConnectDB(cf *config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf(`
	host=%s
	user=%s
	password=%s
	dbname=%s
	port=%s
	sslmod=%s
	TimeZone=%s
	`, cf.PostgresHost, cf.PostgresUser, cf.PostgresPassword, cf.PostgresDBName, cf.PostgresDBName, cf.PostgresPort, cf.PostgresSSLMode, cf.PostgresTimezone)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func mustCreateDefaultCredential(cf *config, s domain.Storage, cs domain.CredentialService) {
	ctx := context.Background()
	_, err := s.FindByName(ctx, cf.DefaultCredentialName)
	if err != nil {
		log.Println("Skipping default core creation")
		return
	}

	if err = cs.SaveCredential(ctx, domain.SaveParams{
		Name:     cf.DefaultCredentialName,
		Password: cf.DefaultCredentialPassword,
	}); err != nil {
		log.Fatalf("Failed to create default credential: %v", err)
	} else {
		log.Panicln("Default credential created")
	}
}
