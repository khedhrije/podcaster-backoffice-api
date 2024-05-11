package mysql

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/khedhrije/podcaster-backoffice-api/internal/configuration"
	"log"
)

// client holds a pointer to a sqlx database connection.
type client struct {
	db *sqlx.DB
}

// NewClient creates a new MySQL client using the provided configuration.
func NewClient(config *configuration.AppConfig) *client {
	db, err := openDB(config.DatabaseConfig.DSN, config.DatabaseConfig)
	if err != nil {
		log.Fatalf("could not open mysql source connection: %s", err.Error())
	}
	return &client{db: db}
}

// openDB opens a new database connection using the provided DSN and database configuration.
func openDB(dsn string, config configuration.DatabaseConfig) (*sqlx.DB, error) {
	conn, err := sqlx.Open(config.Driver, dsn)
	if err != nil {
		return nil, err
	}
	conn.SetMaxOpenConns(config.MaxConnections)
	conn.SetMaxIdleConns(config.MaxIdleConnections)
	conn.SetConnMaxLifetime(config.MaxLifetimeConnections)
	return conn, nil
}
