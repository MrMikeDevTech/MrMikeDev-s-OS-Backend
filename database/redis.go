package database

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	host := getEnv("REDIS_HOST", "127.0.0.1")
	port := getEnv("REDIS_PORT", "15803")
	password := getEnv("REDIS_PASSWORD", "")

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       0,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a Redis: %v", err)
	}

	log.Println("✅ Conexión a Redis establecida en puerto", port)
}
