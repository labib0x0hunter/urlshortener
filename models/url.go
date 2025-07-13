package models

import "time"

// Url represents a shortened URL mapping with metadata.
type Url struct {
	Id        int       // Unique identifier for the URL record
	URL       string    // Original (long) URL
	ShortURL  string    // Generated short code for the URL
	CreatedAt time.Time // Timestamp when the short URL was created
	Expire    time.Time // Expiration time for the short URL (same as CreatedAt if no expiration)
}
