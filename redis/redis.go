package redis

import (
	"context"
	"log/slog"
	"net"
	"os"
	"urlshortener/utils"

	"github.com/redis/go-redis/v9"
)

// NewRedisClient creates and returns a new Redis client connected to the specified database number.
// It reads the host and port from environment variables REDIS_HOST and REDIS_PORT.
// The function pings the Redis server to ensure connectivity and flushes all data before returning the client.
// Returns the Redis client or an error if connection or flush fails.
func NewRedisClient(databaseNo int) (*redis.Client, error) {

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(host, port), // Redis server address
		Password: "",                           // No password set
		DB:       databaseNo,                   // Database number
	})

	// Ping the Redis server to check connectivity
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		slog.Error(" [redis.go] [PING REDIS] ", slog.Any("error", err))
		return nil, utils.ErrDatabaseConnection
	}

	// // Flush all data from the Redis server
	// if err := client.FlushAll(context.Background()).Err(); err != nil {
	// 	slog.Error(" [redis.go] [FLUSH REDIS] ", slog.Any("error", err))
	// 	return nil, utils.ErrDatabaseDelete
	// }

	return client, nil
}
