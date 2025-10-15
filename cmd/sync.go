package cmd

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/zemzale/backscreen-home/storage"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync currency exchange rates",
	Long:  `Sync currency exchange rates from the source to the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement the sync
		ctx := cmd.Context()

		db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/backscreen_home")
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		storage := storage.New(db)
		if err := storage.Migrate(ctx); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}

		return errors.New("not implemented")
	},
}
