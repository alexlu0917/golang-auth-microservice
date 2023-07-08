package main

import (
	"fmt"
	"log"

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
