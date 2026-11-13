package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	"github.com/jayleaton/myPortfolio/packages/database/model"
	"github.com/jayleaton/myPortfolio/packages/database/sqlc"
)

// CacheRepository persists last good live values (cache_entry). It
// implements services.CacheStore.
type CacheRepository struct {
	queries *sqlc.Queries
}

// NewCacheRepository returns a repository backed by db.
func NewCacheRepository(db sqlc.DBTX) *CacheRepository {
	return &CacheRepository{queries: sqlc.New(db)}
}

// GetCache returns model.ErrNotFound when nothing is stored at key.
func (repository *CacheRepository) GetCache(ctx context.Context, key string) (model.CacheEntry, error) {
	entry, err := repository.queries.GetCacheEntry(ctx, key)
	if errors.Is(err, pgx.ErrNoRows) {
		return model.CacheEntry{}, model.ErrNotFound
	}
	if err != nil {
		return model.CacheEntry{}, err
	}
	return model.CacheEntry{Key: entry.Key, Value: entry.Value, UpdatedAt: entry.UpdatedAt}, nil
}

// PutCache stores entry, replacing any value at its key.
func (repository *CacheRepository) PutCache(ctx context.Context, entry model.CacheEntry) error {
	return repository.queries.PutCacheEntry(ctx, sqlc.PutCacheEntryParams{Key: entry.Key, Value: entry.Value, UpdatedAt: entry.UpdatedAt})
}
