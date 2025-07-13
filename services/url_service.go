package services

import (
	"log/slog"
	"time"
	"urlshortener/models"
	"urlshortener/repositories"
	"urlshortener/utils"
)

// UrlService provides methods for URL shortening and retrieval
// It uses a UrlRepository for persistence and lookup
type UrlService struct {
	UrlRepo repositories.UrlRepository // Underlying repository for URL data
}

// NewUrlService creates a new UrlService with the given repository
func NewUrlService(repo repositories.UrlRepository) *UrlService {
	return &UrlService{
		UrlRepo: repo,
	}
}

// CreateShortUrl generates a short URL for the given original URL
// expireIn is the expiration time in minutes (0 means no expiration)
// userAgent is used to help generate a unique short code
// Returns the short code or an error if creation fails
func (u *UrlService) CreateShortUrl(url string, expireIn int64, userAgent string) (string, string, error) {
	uniqueId := utils.UniqueId(userAgent)      // Generate a unique ID based on user agent
	short := utils.GetShortUrl(url + uniqueId) // Generate a short code using the URL and unique ID

	var createdAt, expireAt time.Time
	createdAt = time.Now()
	if expireIn > 0 {
		expireAt = createdAt.Add(time.Duration(expireIn) * time.Minute) // Set expiration if provided
	} else {
		expireAt = createdAt // No expiration, set to creation time
	}

	// Check if the short code already exists (collision check)
	existShort, err := u.UrlRepo.GetByShortCode(short)
	if err != nil {
		slog.Error(" [url_service.go] [CreateShortUrl] ", slog.Any("error", err))
		return "", "", err
	}
	if existShort != nil {
		return "", "", utils.ErrShortCodeCollision
	}

	// Create the Url model
	shortUrl := models.Url{
		URL:       url,
		ShortURL:  short,
		CreatedAt: createdAt,
		Expire:    expireAt,
	}

	err = u.UrlRepo.Create(shortUrl)
	if err != nil {
		slog.Error(" [url_service.go] [CREATE] ", slog.Any("error", err))
	}

	expireMsg := "no expiration"
	if expireIn > 0 {
		expireMsg = expireAt.Format("2006-01-02 15:04:05") // Format expiration time if set
	}

	return short, expireMsg, err
}

// GetUrlByCode retrieves the original URL by its short code
// Returns the Url model or an error if not found
func (u *UrlService) GetUrlByCode(code string) (*models.Url, error) {
	url, err := u.UrlRepo.GetByShortCode(code)
	if err != nil {
		slog.Error(" [url_service.go] [GetUrlByCode] ", slog.Any("error", err))
		return nil, err
	}
	if url == nil {
		return nil, utils.ErrUrlNotFound
	}
	return url, nil
}
