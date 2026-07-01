package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

var ClientRedis *redis.Client

func ConectarRedis() error {

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		return fmt.Errorf("a variável de ambiente REDIS_URL não está configurada")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		panic(err)
	}

	ClientRedis = redis.NewClient(opt)

	ping, err := ClientRedis.Ping(context.Background()).Result()

	if err != nil {
		return fmt.Errorf("falha ao conectar ao Redis: %v", err)
	}

	fmt.Println("Conexão com o Redis estabelecida com sucesso!" + ping)
	return nil
}

func FecharRedis(){
	if ClientRedis != nil {
		ClientRedis.Close()
	}
}