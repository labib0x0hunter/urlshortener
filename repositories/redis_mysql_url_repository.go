package repositories

import (
	"encoding/json"
	"log/slog"
	"time"
	"urlshortener/cache"
	"urlshortener/models"
	"urlshortener/utils"
)

// RedisMysqlUrlRepository is a URL repository that uses both a persistent backend (MySQL) and a cache (Redis)
type RedisMysqlUrlRepository struct {
	repo  UrlRepository // Underlying persistent repository (e.g., MySQL)
	redis cache.Cache   // Cache layer (e.g., Redis)
}

// NewRedisMysqlUrlRepository creates a new RedisMysqlUrlRepository with the given persistent repo and cache
func NewRedisMysqlUrlRepository(repo UrlRepository, redis cache.Cache) *RedisMysqlUrlRepository {
	return &RedisMysqlUrlRepository{
		repo:  repo,
		redis: redis,
	}
}

// Create stores a new URL mapping in the persistent repository
func (r *RedisMysqlUrlRepository) Create(url models.Url) error {
	return r.repo.Create(url)
}

// GetByShortCode retrieves a URL by its short code, using cache and expiration logic
func (r *RedisMysqlUrlRepository) GetByShortCode(shortCode string) (*models.Url, error) {
	// Check if the short code is marked as expired in cache
	expireKey := "expire:" + shortCode
	_, err := r.redis.Get(expireKey)
	if err == nil {
		slog.Error(" [redis_mysql_url_repository.go] [Expired short code] ", slog.Any("error", err))
		return nil, utils.ErrShortCodeExpired
	}

	// Try to get the URL from cache
	cacheKey := "short:" + shortCode
	value, err := r.redis.Get(cacheKey)
	if err == nil {
		var cacheUrl models.Url
		if err := json.Unmarshal([]byte(value), &cacheUrl); err == nil {
			return &cacheUrl, nil // Return cached URL if found
		}
	}

	// slog.Info(" [Cache miss] ", slog.String("shortCode", shortCode))

	// Fallback to persistent repository if not in cache
	url, err := r.repo.GetByShortCode(shortCode)
	if err != nil || url == nil {
		return nil, err
	}

	// slog.Info(" [URL found in repo] ", slog.String("shortCode", shortCode))

	// If the URL is expired, mark it as expired in cache and return error
	if url.CreatedAt != url.Expire && url.Expire.Before(time.Now()) {
		expireKey := "expire:" + shortCode
		r.redis.Set(expireKey, "1", 0) // 0 means never expire in cache
		return nil, utils.ErrShortCodeExpired
	}

	// Cache the URL for future lookups min(5 minutes, actual expiration)
	urlJson, err := json.Marshal(url)
	if err == nil {

		// slog.Info(" [Created At]", url.CreatedAt.GoString())
		// slog.Info(" [Expire At]", url.Expire.GoString())

		duration := 5 * time.Minute
		if url.CreatedAt != url.Expire {
			expire := url.Expire.Sub(url.CreatedAt)
			duration = min(expire, duration)
		}

		// slog.Info(" [Caching URL] ", slog.String("shortCode", shortCode), slog.String("duration", duration.String()))

		r.redis.Set(cacheKey, string(urlJson), duration)
	}

	return url, err
}
