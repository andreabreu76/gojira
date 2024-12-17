package main

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

func InitializeRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	fmt.Println("[INFO] Conexão com Redis estabelecida.")
}

func StoreResult(key string, value string) error {
	return rdb.Set(ctx, key, value, 0).Err()
}

func FetchAllResults(prefix string) map[string]string {
	keys, _ := rdb.Keys(ctx, prefix+"*").Result()
	results := make(map[string]string)

	for _, key := range keys {
		value, _ := rdb.Get(ctx, key).Result()
		results[key] = value
	}

	return results
}
