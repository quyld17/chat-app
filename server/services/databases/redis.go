package databases

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func NewRedis() *redis.Client {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Failed to retrieve credentials. Please try again")
	}

	dbRedis := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT_REDIS"),
		Password: "",
		DB:       0,
	})

	return dbRedis
}
