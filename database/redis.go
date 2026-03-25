package database

import (
	"context"
	"fmt"
	"log"

	"github.com/MrMikeDevTech/mrmikedevs-os/utils"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func ConnectRedis() {
	host := utils.GetEnv("REDIS_HOST", "redis")
	port := utils.GetEnv("REDIS_PORT", "6379")
	password := utils.GetEnv("REDIS_PASSWORD")

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
