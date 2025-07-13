package handlers

import (
	"net/url"
	"os"
	"urlshortener/services"
	"urlshortener/utils"

	"github.com/gin-gonic/gin"
)

// ShortenHandler handles HTTP requests for URL shortening and redirection
// It uses a UrlService to perform business logic
type ShortenHandler struct {
	UrlService *services.UrlService // Service for URL operations
}

// UrlRequest represents the expected JSON payload for shortening a URL
// ExpireAt is optional and specifies expiration in minutes
type UrlRequest struct {
	Url      string `json:"url" binding:"required"` // The original URL to shorten
	ExpireAt int64  `json:"expire_in,omitempty"`    // Expiration in minutes (optional)
}

// NewShortenHandler creates a new ShortenHandler with the given UrlService
func NewShortenHandler(UrlService *services.UrlService) *ShortenHandler {
	return &ShortenHandler{
		UrlService: UrlService,
	}
}

// ShortenURL handles POST /shorten requests to create a new short URL
// Validates input, calls the service, and returns the result as JSON
func (s *ShortenHandler) ShortenURL(ctx *gin.Context) {
	var req UrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(422, gin.H{
			"error": "Validation error",
		})
		return
	}

	// Validate URL format
	u, err := url.ParseRequestURI(req.Url)
	if err != nil || u.Scheme == "" || u.Host == "" {
		ctx.JSON(400, gin.H{
			"error": "Invalid URL format",
		})
		return
	}

	// Check if URL length is valid (greater than 25 characters)
	if len(req.Url) <= 25 {
		ctx.JSON(400, gin.H{
			"error": "URL must be longer than 25 characters",
		})
		return
	}

	// Create the short URL using the service
	short, expireAt, err := s.UrlService.CreateShortUrl(req.Url, req.ExpireAt, ctx.Request.UserAgent())
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Internal server error",
		})
		return
	}

	// Build the full short URL with prefix
	prefix := os.Getenv("SHORT_URL_PREFIX")
	ctx.JSON(201, gin.H{
		"message":   "success",
		"short_url": prefix + short,
		"expire_at": expireAt,
	})
}

// GetFullURL handles GET /:code requests to redirect to the original URL
// Looks up the short code and redirects, or returns an error if not found
func (s *ShortenHandler) GetFullURL(ctx *gin.Context) {
	shortCode := ctx.Param("code")
	if shortCode == "" {
		ctx.JSON(400, gin.H{
			"error": "Short URL code is required",
		})
		return
	}
	url, err := s.UrlService.GetUrlByCode(shortCode)

	if err != nil {
		errMsg := "Internal server error"
		errCode := 500
		switch err {
		case utils.ErrUrlNotFound:
			errMsg, errCode = "URL not found", 404

		case utils.ErrShortCodeExpired:
			errMsg, errCode = "URL has expired", 410
		}

		ctx.JSON(errCode, gin.H{
			"error": errMsg,
		})
		return
	}

	ctx.Redirect(301, url.URL)
}

// GetUrlMetadata handles GET /fetch/:code requests to retrieve URL metadata
// Looks up the short code and returns metadata without redirecting
func (s *ShortenHandler) GetUrlMetadata(ctx *gin.Context) {
	shortCode := ctx.Param("code")
	if shortCode == "" {
		ctx.JSON(400, gin.H{
			"error": "Short URL code is required",
		})
		return
	}
	url, err := s.UrlService.GetUrlByCode(shortCode)

	if err != nil && err != utils.ErrUrlNotFound {
		errMsg := "Internal server error"
		errCode := 500
		switch err {
		case utils.ErrUrlNotFound:
			errMsg, errCode = "URL not found", 404

		// case utils.ErrShortCodeExpired:
		// 	errMsg, errCode = "URL has expired", 410
		}

		ctx.JSON(errCode, gin.H{
			"error": errMsg,
		})
		return
	}

	ctx.JSON(200, gin.H{
		"url": url.URL,
		"metadata": gin.H{
			"short_code": shortCode,
			"created_at": url.CreatedAt,
			"expire_at":  url.Expire,
		},
	})
}