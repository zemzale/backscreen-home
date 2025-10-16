package syncer

import (
	"context"
	"errors"
	"log/slog"
	"sync"

	"github.com/zemzale/backscreen-home/domain/entity"
	"github.com/zemzale/backscreen-home/storage"
)

type RateFetcher interface {
	Fetch(ctx context.Context, currency string) ([]entity.Rate, error)
}

type Usecase struct {
	store   *storage.Client
	fetcher RateFetcher
}

func New(store *storage.Client, fetcher RateFetcher) *Usecase {
	return &Usecase{
		store:   store,
		fetcher: fetcher,
	}
}

func (u *Usecase) Sync(ctx context.Context, currencies []string) {
	logger := slog.With(slog.String("component", "sync"))

	// TODO: Remove the WaitGroup and use some channels man
	wg := sync.WaitGroup{}
	wg.Add(len(currencies))

	for _, currency := range currencies {
		go func(wg *sync.WaitGroup, currency string) {
			defer wg.Done()

			logger.InfoContext(ctx, "Syncing currency", slog.String("currency", currency))

			u.syncCurrency(ctx, currency)
		}(&wg, currency)
	}

	wg.Wait()
}

func (u *Usecase) syncCurrency(ctx context.Context, currency string) {
	logger := slog.With(slog.String("component", "sync"), slog.String("currency", currency))

	rates, err := u.fetcher.Fetch(ctx, currency)
	if err != nil {
		logger.ErrorContext(ctx, "Failed to fetch rate", slog.Any("currency", currency), slog.Any("error", err))
		return
	}

	// TODO: Batch the rate inserts, to get more performance
	logger.DebugContext(ctx, "Storing rates to database", slog.Any("rates", rates))
	for _, rate := range rates {
		if err := u.store.StoreRate(ctx, rate); err != nil {
			if errors.Is(err, storage.ErrDuplicate) {
				logger.WarnContext(ctx,
					"Rate already exists in database",
					slog.Any("rate", rate),
				)
				continue
			}
			logger.ErrorContext(ctx, "Failed to store rate", slog.Any("rate", rate), slog.Any("error", err))
		}
	}
}
