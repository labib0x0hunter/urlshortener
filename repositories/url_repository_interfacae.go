package repositories

import "urlshortener/models"

// UrlRepository defines the interface for URL persistence and retrieval.
// Implementations may use different storage backends (e.g., MySQL, Redis, etc.).
type UrlRepository interface {
	// Create stores a new URL mapping in the repository.
	Create(url models.Url) error
	// GetByShortCode retrieves a URL mapping by its short code.
	GetByShortCode(shortCode string) (*models.Url, error)
}
