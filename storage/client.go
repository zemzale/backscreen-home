package storage

import (
	"context"
	"errors"
	"log/slog"
	"time"

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

type Rate struct {
	ID          int       `db:"id"`
	Code        string    `db:"code"`
	Value       string    `db:"value"`
	PunlishedAt time.Time `db:"published_at"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func (r Rate) ToEntity() entity.Rate {
	return entity.Rate{
		Code:        r.Code,
		Value:       r.Value,
		PublishedAt: r.PunlishedAt,
	}
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

func (c *Client) GetLatestRate(ctx context.Context, code string) (entity.Rate, error) {
	var rate Rate

	err := c.db.GetContext(ctx, &rate, `
		SELECT code, value, published_at FROM rates WHERE code = ? ORDER BY published_at DESC LIMIT 1;
	`, code)
	if err != nil {
		return entity.Rate{}, err
	}

	return rate.ToEntity(), nil
}
