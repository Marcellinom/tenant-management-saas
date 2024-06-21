package main

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal"
	"github.com/Marcellinom/tenant-management-saas/provider"
	"github.com/Marcellinom/tenant-management-saas/provider/auth"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("tidak ada .env file yang terdeteksi", err)
	}
	startApp()
}

func startApp() {
	engine_cfg := provider.DefaultEngineConfig()
	engine_cfg.UseCustomErrorHandler(errors.DefaultHandler())
	engine, err := provider.SetupWebEngine(engine_cfg)
	if err != nil {
		log.Panic(err)
	}

	tm_database := provider.NewConnectionConfig(
		os.Getenv("DB_CONNECTION"),
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	iam_database := provider.NewConnectionConfig(
		os.Getenv("IAM_DB_CONNECTION"),
		os.Getenv("IAM_DB_DRIVER"),
		os.Getenv("IAM_DB_USER"),
		os.Getenv("IAM_DB_PASSWORD"),
		os.Getenv("IAM_DB_HOST"),
		os.Getenv("IAM_DB_PORT"),
		os.Getenv("IAM_DB_DATABASE"),
	)

	db_connections, err := provider.SetupDatabase([]provider.ConnectionConfig{
		tm_database,
		iam_database,
	})
	if err != nil {
		log.Panic(err)
	}

	iam, err := auth.New(os.Getenv("AUTH_PROVIDER"))
	if err != nil {
		log.Panic(err)
	}

	app := provider.NewApplication(engine, db_connections, iam)

	internal.RegisterApplication(app)

	if engine_cfg.Port == "" {
		engine_cfg.Port = "8080"
	}
	engine.Run(":" + engine_cfg.Port)
}
