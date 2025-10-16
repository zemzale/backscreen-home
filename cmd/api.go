package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog/v3"
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
		mux.Use(httplog.RequestLogger(logger, &httplog.Options{
			Level:         slog.LevelDebug,
			Schema:        httplog.SchemaECS,
			RecoverPanics: true,
		}))
		handler := server.HandlerFromMux(
			server.NewStrictHandler(api{store: store}, nil),
			mux,
		)

		errChan := make(chan error, 1)
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		server := &http.Server{
			Addr:    host,
			Handler: handler,
		}

		go func() {
			errChan <- server.ListenAndServe()
			close(errChan)
		}()

		select {
		case err := <-errChan:
			if err != nil {
				if !errors.Is(err, http.ErrServerClosed) {
					logger.ErrorContext(ctx, "API server closed with error", slog.String("err", err.Error()))
				}
				slog.InfoContext(ctx, "API stoped serving new conntections")
			}
		case <-ctx.Done():
			return fmt.Errorf("context done: %w", ctx.Err())
		case <-sigChan:
			logger.InfoContext(ctx, "API received SIGINT or SIGTERM")

			shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()

			if err := server.Shutdown(shutdownCtx); err != nil {
				logger.ErrorContext(ctx, "Failed to shutdown API", slog.String("err", err.Error()))
			}

			slog.InfoContext(ctx, "API shutdown")
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
		if errors.Is(err, storage.ErrNotFound) {
			return server.GetApiV1Currency404Response{}, nil
		}

		return server.GetApiV1Currency500JSONResponse{
			InternalServerErrorJSONResponse: errToInternalServerError(err),
		}, nil
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
		if errors.Is(err, storage.ErrNotFound) {
			return server.GetApiV1Currency404Response{}, nil
		}

		return server.GetApiV1CurrencyHistory500JSONResponse{
			InternalServerErrorJSONResponse: errToInternalServerError(err),
		}, nil
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

func errToInternalServerError(err error) server.InternalServerErrorJSONResponse {
	errStr := err.Error()
	return server.InternalServerErrorJSONResponse{
		Error: &errStr,
	}
}
