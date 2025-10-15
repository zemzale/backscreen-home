package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/zemzale/backscreen-home/domain/entity"
)

const createRatesTable = `
CREATE TABLE IF NOT EXISTS rates (
	id INT NOT NULL AUTO_INCREMENT,
	code VARCHAR(3) NOT NULL,
	value VARCHAR(100) NOT NULL,
	published_at DATETIME NOT NULL,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	PRIMARY KEY (id),
	UNIQUE INDEX (code, published_at),
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

func (c *Client) StoreRate(ctx context.Context, rate entity.Rate) error {
	_, err := c.db.ExecContext(ctx, `
		INSERT INTO rates (code, value, published_at) VALUES (?, ?, ?);
	`, rate.Code, rate.Value, rate.PublishedAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 {
				return ErrDuplicate
			}
		}

	}

	return err
}
