package main

import (
	"urlshortener/handlers"

	"github.com/gin-gonic/gin"
)

func main() {


	shortenService := ser
	shortenHandler := handlers.NewShortenHandler()

	router := gin.Default()

	router.GET("/:code", shortenHandler.GetFullURL)
	router.POST("/shorten", shortenHandler.ShortenURL)

	router.Run(":3000")

}
