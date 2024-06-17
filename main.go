package main

import (
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
		log.Panic("Error loading .env file", err)
	}
	//gcp.Publish(context.TODO())
	//go gcp.Subscribe(context.TODO())
	//testTerraform()
	startApp()
}

func testTerraform() {
	//tenant_id := "alex"
	////gcs_bucket := gcp.Bucket("saas-tf-config", "tenants")
	//terraform.New("E:\\1kuliah\\TA\\code\\tenant-management\\terraform").
	//	NewTenantConfig(
	//		tenant_id,
	//		"list-foto-product",
	//		terraform.SILO,
	//		*tfexec.Var(fmt.Sprintf("tenant_id=%s", tenant_id)),
	//	).
	//	UseBackend(terraform.BuiltinBackend("saas-tf-config", "tenants")).
	//	//UseBackend(gcs_bucket).
	//	Apply()
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

	db_connections, err := provider.SetupDatabase([]provider.ConnectionConfig{
		tm_database,
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
