package cmd

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/zemzale/backscreen-home/domain/usecase/syncer"
	"github.com/zemzale/backscreen-home/sources/lvbank"
)

var allowedCurrencies = []string{"AUD", "BGN", "BRL", "CAD", "CHF", "CNY", "CZK", "DKK", "GBP", "HKD"}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync currency exchange rates",
	Long:  `Sync currency exchange rates from the source to the database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: Implement the sync
		ctx := cmd.Context()

		logger := slog.With("component", "sync")

		logger.InfoContext(ctx, "Starting syncing currencies")

		syncer.New(store, lvbank.New()).Sync(ctx, allowedCurrencies)
		logger.InfoContext(ctx, "Finished syncing currencies")

		return nil
	},
}
