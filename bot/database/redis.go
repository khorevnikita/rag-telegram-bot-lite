package database

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"log"
)

var Redis *redis.Client

func RedisConnect() {
	Redis = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
	ctx := context.Background()
	if err := Redis.Ping(ctx).Err(); err != nil {
		log.Fatalf("Не удалось подключиться к Redis: %v", err)
	}
}

func Publish(channel string, payload interface{}) error {
	ctx := context.Background()
	message, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return Redis.Publish(ctx, channel, message).Err()
}

func Subscribe[T any](ctx context.Context, channel string, cb func(T)) {
	subscriber := Redis.Subscribe(ctx, channel)
	ch := subscriber.Channel()

	log.Printf("Подписываемся на %s\n", channel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("Завершаем подписку на %s\n", channel)
			return
		case msg := <-ch:
			go func(payload string) {
				var body T
				err := json.Unmarshal([]byte(payload), &body)
				if err != nil {
					log.Printf("Ошибка десериализации: %v", err)
					return
				}
				cb(body)
			}(msg.Payload)
		}
	}
}

func CloseRedis() {
	if err := Redis.Close(); err != nil {
		log.Printf("Ошибка при закрытии Redis: %v", err)
	}
}
