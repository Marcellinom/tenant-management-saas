package provider

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database = gorm.DB

type Connection map[string]*Database

type ConnectionConfig struct {
	ConnectionName string
	Driver         string
	User           string
	Password       string
	Host           string
	Port           string
	Database       string
}

func NewConnectionConfig(connectionName string, driver string, user string, password string, host string, port string, database string) ConnectionConfig {
	return ConnectionConfig{ConnectionName: connectionName, Driver: driver, User: user, Password: password, Host: host, Port: port, Database: database}
}

func SetupDatabase(configs []ConnectionConfig) (*Connection, error) {
	connection := make(Connection)
	for _, cfg := range configs {
		if cfg.Driver == "" {
			return nil, fmt.Errorf("database driver is empty, supported drivers are [postgres]")
		}
		switch cfg.Driver {
		case "sqlite":
			return nil, fmt.Errorf("database driver %s is inder development, available=[postgres]", cfg.Driver)
		case "sqlserver":
			return nil, fmt.Errorf("database driver %s is inder development, available=[postgres]", cfg.Driver)
		case "postgres":
			dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
			// log.Println("Connecting to PostgreSQL database...")
			db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err != nil {
				return nil, fmt.Errorf("PostgreSQL connection error: %w", err)
			}
			connection[cfg.ConnectionName] = db
		default:
			return nil, fmt.Errorf("unknown database driver %s, supported drivers are [sqlite, sqlserver, postgres]", cfg.Driver)
		}
	}
	return &connection, nil
}
