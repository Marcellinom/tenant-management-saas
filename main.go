package main

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"tenant_management/internal"
	"tenant_management/pkg"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file")
	}

	//tenant_id := "marsel"
	//gcs_bucket := gcp.Bucket("saas-tf-config", "tenants")
	//terraform.New("E:\\1kuliah\\TA\\code\\tenant-management\\terraform").
	//	Tenant(
	//		tenant_id,
	//		"list-foto-product",
	//		terraform.SILO,
	//		tfexec.Var("tenant_id=marsel")).
	//	UseBackend(gcs_bucket).
	//	Create()
	//
	//os.Exit(1)

	engine_cfg := pkg.DefaultEngineConfig()
	engine, err := pkg.SetupWebEngine(engine_cfg)
	if err != nil {
		log.Panic(err)
	}

	pool_config := pkg.NewConnectionConfig(
		os.Getenv("DB_CONNECTION"),
		os.Getenv("DB_DRIVER"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)

	db_connections, err := pkg.SetupDatabase([]pkg.ConnectionConfig{
		pool_config,
	})
	if err != nil {
		log.Panic(err)
	}
	app := pkg.NewApplication(engine, db_connections)
	internal.RegisterApplication(app)

	if engine_cfg.Port == "" {
		engine_cfg.Port = "8080"
	}
	engine.Run(":" + engine_cfg.Port)
}
