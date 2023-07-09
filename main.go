package main

import (
	"context"
	"fmt"
	"log"
	"microauth/domain"
	"microauth/infra"
	"microauth/rest"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cf := loadConfig()
	db := mustConnectDB(cf)
	ps, err := infra.NewPostgressStorage(db)
	if err != nil {
		log.Fatalf("Failed to create postgres: %v", err)
	}
	var bh domain.Hasher = infra.BCryptHasher{}

	authSrv := domain.NewAuthService(ps, bh)
	credentialSrc := domain.NewCredentialService(ps, bh, authSrv)

	mustCreateDefaultCredential(cf, ps, credentialSrc)
	authHr := rest.NewAuthHandler(authSrv)
	authMd := rest.NewAuthMiddleware(authSrv)

	e := echo.New()

	api := e.Group("/api/v1")

	api.POST("/login", authHr.HandleLogin)
	api.POST("/logout", authHr.HandleLogout)

	dash := api.Group("/dashboard")
	dash.Use(authMd)

	e.Logger.Fatal(e.Start(":9876"))
}

func mustConnectDB(cf *config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s	password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", cf.PostgresHost, cf.PostgresUser, cf.PostgresPassword, cf.PostgresDBName, cf.PostgresPort)), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func mustCreateDefaultCredential(cf *config, s domain.Storage, cs domain.CredentialService) {
	ctx := context.Background()
	_, err := s.FindByName(ctx, cf.DefaultCredentialName)
	if err == nil {
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
