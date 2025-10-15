package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"github.com/zemzale/backscreen-home/domain/entity"
	"github.com/zemzale/backscreen-home/domain/mapper"
	"github.com/zemzale/backscreen-home/storage"
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

		logger.DebugContext(ctx, "Connecting to database")
		db, err := sqlx.Connect("mysql", "root:root@tcp(localhost:3306)/backscreen_home")
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}

		logger.DebugContext(ctx, "Creating storage client")
		storage := storage.New(db)
		if err := storage.Migrate(ctx); err != nil {
			return fmt.Errorf("failed to migrate database: %w", err)
		}

		logger.InfoContext(ctx, "Starting syncing currencies")

		wg := sync.WaitGroup{}
		wg.Add(len(allowedCurrencies))

		for _, currency := range allowedCurrencies {
			go func(wg *sync.WaitGroup, currency string) {
				defer wg.Done()

				logger.InfoContext(ctx, "Syncing currency", slog.String("currency", currency))

				if err := syncCurrency(ctx, storage, currency, LVBankRSSRateFetcher{}); err != nil {
					logger.ErrorContext(ctx, "Failed to sync currency", slog.String("currency", currency), slog.Any("error", err))
				}
			}(&wg, currency)
		}

		wg.Wait()

		logger.InfoContext(ctx, "Finished syncing currencies")

		return nil
	},
}

type RateFetcher interface {
	Fetch(ctx context.Context, currency string) (entity.Rate, error)
}

func syncCurrency(ctx context.Context, store *storage.Client, currency string, rateFetcher RateFetcher) error {
	logger := slog.With(slog.String("component", "sync"), slog.String("currency", currency))

	rate, err := rateFetcher.Fetch(ctx, currency)
	if err != nil {
		return fmt.Errorf("failed to fetch rate: %w", err)
	}

	logger.DebugContext(ctx, "Storing rates to database", slog.Any("rate", rate))
	if err := store.StoreRate(ctx, rate); err != nil {
		if errors.Is(err, storage.ErrDuplicate) {
			logger.InfoContext(ctx,
				"Rate already exists in database",
				slog.Any("rate", rate),
			)
			return nil
		}
		return fmt.Errorf("failed to store rate: %w", err)
	}

	return nil
}

// TODO: Maybe just move this to some other package to make it clearer
type LVBankRSSRateFetcher struct{}

func (f LVBankRSSRateFetcher) Fetch(ctx context.Context, currency string) (entity.Rate, error) {
	logger := slog.With("component", "LVBankRSSRateFetcher", "currency", currency)
	logger.DebugContext(ctx, "Fetching rates")

	url := "https://www.bank.lv/vk/ecb_rss.xml"
	logger.DebugContext(ctx, "Creating requets for exchange rates", slog.String("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return entity.Rate{}, fmt.Errorf("failed to create request: %w", err)
	}

	// TODO: Change out the default client, since it's very bad and has no timeouts

	logger.InfoContext(ctx, "Sending request for exchange rates", slog.String("url", url))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return entity.Rate{}, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	logger.InfoContext(
		ctx,
		"Received response for exchange rates",
		slog.String("url", url),
		slog.Int("status", resp.StatusCode),
	)
	if resp.StatusCode != http.StatusOK {
		return entity.Rate{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	logger.DebugContext(ctx, "Parsing rates")
	rates, err := mapper.RatesFromXML(resp.Body)
	if err != nil {
		return entity.Rate{}, fmt.Errorf("failed to parse rates: %w", err)
	}

	logger.DebugContext(ctx, "Searchign for rate in response", slog.Int("rate_count", len(rates)))
	idx := slices.IndexFunc(rates, func(r entity.Rate) bool {
		return r.Code == currency
	})

	logger.DebugContext(ctx, "Found rate", slog.Int("index", idx))
	if idx == -1 {
		return entity.Rate{}, fmt.Errorf("currency %s not found", currency)
	}

	return rates[idx], nil
}
