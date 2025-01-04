package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func getString(key, fallback string) string {
	val, err := os.LookupEnv(key)
	if err {
		return fallback
	}
	return val
}

func getInt(key string, fallback int) int {
	val, err := os.LookupEnv(key)
	if err {
		return fallback
	}
	valAsInt, er := strconv.Atoi(val)
	if er != nil {
		return fallback
	}

	return valAsInt
}
