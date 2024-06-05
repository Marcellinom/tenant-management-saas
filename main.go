package main

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/internal"
	"github.com/Marcellinom/tenant-management-saas/pkg"
	"github.com/Marcellinom/tenant-management-saas/pkg/auth"
	"github.com/Marcellinom/tenant-management-saas/pkg/terraform"
	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/joho/godotenv"
	"log"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Panic("Error loading .env file", err)
	}
	//startApp()
	testTerraform()
}

func testTerraform() {
	tenant_id := "alex"
	//gcs_bucket := gcp.Bucket("saas-tf-config", "tenants")
	terraform.New("E:\\1kuliah\\TA\\code\\tenant-management\\terraform").
		Tenant(
			tenant_id,
			"list-foto-product",
			terraform.SILO,
			*tfexec.Var(fmt.Sprintf("tenant_id=%s", tenant_id)),
		).
		UseBackend(terraform.BuiltinBackend("saas-tf-config", "tenants")).
		//UseBackend(gcs_bucket).
		Apply()
}

func startApp() {
	engine_cfg := pkg.DefaultEngineConfig()
	engine, err := pkg.SetupWebEngine(engine_cfg)
	if err != nil {
		log.Panic(err)
	}

	//pool_config := pkg.NewConnectionConfig(
	//	os.Getenv("DB_CONNECTION"),
	//	os.Getenv("DB_DRIVER"),
	//	os.Getenv("DB_USER"),
	//	os.Getenv("DB_PASSWORD"),
	//	os.Getenv("DB_HOST"),
	//	os.Getenv("DB_PORT"),
	//	os.Getenv("DB_DATABASE"),
	//)
	//
	//db_connections, err := pkg.SetupDatabase([]pkg.ConnectionConfig{
	//	pool_config,
	//})
	if err != nil {
		log.Panic(err)
	}
	app := pkg.NewApplication(engine, nil)

	iam, err := auth.New()
	if err != nil {
		log.Panic(err)
	}

	internal.RegisterApplication(app)
	iam.RegisterCallback(engine)

	if engine_cfg.Port == "" {
		engine_cfg.Port = "80"
	}
	engine.Run(":" + engine_cfg.Port)
}
