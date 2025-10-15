package cmd

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/zemzale/backscreen-home/pkg/server"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API",
	Long:  `Start the API to get stored currency exchange rates.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		mux := chi.NewRouter()
		handler := server.HandlerFromMux(server.NewStrictHandler(api{}, nil), mux)

		if err := http.ListenAndServe(":8080", handler); err != nil {
			return err
		}

		return nil
	},
}

var _ server.StrictServerInterface = &api{}

type api struct{}

// Get latest exchange rate
// (GET /api/v1/{currency})
func (a api) GetApiV1Currency(ctx context.Context, req server.GetApiV1CurrencyRequestObject) (server.GetApiV1CurrencyResponseObject, error) {
	return server.GetApiV1Currency200JSONResponse{}, nil
}

// Get all historical exchange rates
// (GET /api/v1/{currency}/history)
func (a api) GetApiV1CurrencyHistory(ctx context.Context, req server.GetApiV1CurrencyHistoryRequestObject) (server.GetApiV1CurrencyHistoryResponseObject, error) {
	return server.GetApiV1CurrencyHistory200JSONResponse{}, nil
}
