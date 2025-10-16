package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/zemzale/backscreen-home/domain/entity"
	"github.com/zemzale/backscreen-home/domain/mapper"
	"github.com/zemzale/backscreen-home/domain/usecase/syncer"
	"github.com/zemzale/backscreen-home/slices"
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

		syncer.New(store, LVBankRSSRateFetcher{}).Sync(ctx, allowedCurrencies)
		logger.InfoContext(ctx, "Finished syncing currencies")

		return nil
	},
}

type RateFetcher interface {
	Fetch(ctx context.Context, currency string) ([]entity.Rate, error) // TODO This shoudl return a list of rates, since one rquest can return multiple days
}

// TODO: Maybe just move this to some other package to make it clearer
type LVBankRSSRateFetcher struct{}

func (f LVBankRSSRateFetcher) Fetch(ctx context.Context, currency string) ([]entity.Rate, error) {
	logger := slog.With("component", "LVBankRSSRateFetcher", "currency", currency)
	logger.DebugContext(ctx, "Fetching rates")

	url := "https://www.bank.lv/vk/ecb_rss.xml"
	logger.DebugContext(ctx, "Creating requets for exchange rates", slog.String("url", url))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// TODO: Change out the default client, since it's very bad and has no timeouts

	logger.InfoContext(ctx, "Sending request for exchange rates", slog.String("url", url))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch rates: %w", err)
	}
	defer resp.Body.Close()

	logger.InfoContext(
		ctx,
		"Received response for exchange rates",
		slog.String("url", url),
		slog.Int("status", resp.StatusCode),
	)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	logger.DebugContext(ctx, "Parsing rates")
	rates, err := mapper.RatesFromXML(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rates: %w", err)
	}

	logger.DebugContext(ctx, "Searchign for rate in response", slog.Int("rate_count", len(rates)))

	return slices.FilterInPlace(rates, func(r entity.Rate) bool {
		return r.Code == currency
	}), nil
}
