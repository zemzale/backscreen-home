package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/zemzale/backscreen-home/domain/entity"
	"github.com/zemzale/backscreen-home/pkg/server"
	"github.com/zemzale/backscreen-home/slices"
	"github.com/zemzale/backscreen-home/storage"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API",
	Long:  `Start the API to get stored currency exchange rates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		logger := slog.With("component", "api")

		host := viper.GetString("api.host")

		logger.InfoContext(ctx, "Starting API", slog.String("host", host))

		mux := chi.NewRouter()
		handler := server.HandlerFromMux(server.NewStrictHandler(api{store: store}, nil), mux)

		// TODO Move the server to a gorutine
		// TODO Add a graceful shutdown
		// TODO Add logging
		if err := http.ListenAndServe(host, handler); err != nil {
			return err
		}

		return nil
	},
}

var _ server.StrictServerInterface = &api{}

type api struct {
	store *storage.Client
}

// Get latest exchange rate
// (GET /api/v1/{currency})
func (a api) GetApiV1Currency(ctx context.Context, req server.GetApiV1CurrencyRequestObject) (server.GetApiV1CurrencyResponseObject, error) {
	rate, err := a.store.GetLatestRate(ctx, req.Currency)
	if err != nil {
		return server.GetApiV1Currency200JSONResponse{}, fmt.Errorf("failed to get rate: %w", err)
	}

	return server.GetApiV1Currency200JSONResponse{
		Code:        rate.Code,
		PublishedAt: rate.PublishedAt,
		Value:       rate.Value,
	}, nil
}

// Get all historical exchange rates
// (GET /api/v1/{currency}/history)
func (a api) GetApiV1CurrencyHistory(ctx context.Context, req server.GetApiV1CurrencyHistoryRequestObject) (server.GetApiV1CurrencyHistoryResponseObject, error) {
	rates, err := a.store.GetRates(ctx, req.Currency)
	if err != nil {
		return server.GetApiV1CurrencyHistory200JSONResponse{}, fmt.Errorf("failed to get rates: %w", err)
	}

	ratesResponse := slices.Map(rates, mapRateToApiV1CurrencyHistoryRate)

	return server.GetApiV1CurrencyHistory200JSONResponse(ratesResponse), nil
}

func mapRateToApiV1CurrencyHistoryRate(rate entity.Rate) server.Rate {
	return server.Rate{
		Code:        rate.Code,
		Value:       rate.Value,
		PublishedAt: rate.PublishedAt,
	}
}
