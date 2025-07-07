package handlers

import "github.com/gin-gonic/gin"

type ShortenHandler struct {
	// Service
}

type ShortenUrlRequest struct {
	Url string `json:"url" binding:"required"`
}

func NewShortenHandler() *ShortenHandler {
	return &ShortenHandler{}
}

func (s *ShortenHandler) ShortenURL(ctx *gin.Context) {

}

func (s *ShortenHandler) GetFullURL(ctx *gin.Context) {

}
