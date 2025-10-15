package storage

import (
	"context"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const createRatesTable = `
CREATE TABLE IF NOT EXISTS rates (
	id INT NOT NULL AUTO_INCREMENT,
	code VARCHAR(3) NOT NULL,
	value VARCHAR(10) NOT NULL,
	published_at DATETIME NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	PRIMARY KEY (id),
	INDEX (code),
	INDEX (published_at)
);`

const createImportDataTable = `
CREATE TABLE IF NOT EXISTS import_data (
	id INT NOT NULL AUTO_INCREMENT,
	data TEXT NOT NULL,
	source VARCHAR(255) NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	PRIMARY KEY (id)
);
`

var schemas = []string{
	createRatesTable,
	createImportDataTable,
}

type Client struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) Migrate(ctx context.Context) error {
	logger := slog.With("component", "db")
	logger.DebugContext(ctx, "Running DB migrations")

	for _, query := range schemas {
		_, err := c.db.ExecContext(ctx, query)
		if err != nil {
			return err
		}
	}

	return nil
}
