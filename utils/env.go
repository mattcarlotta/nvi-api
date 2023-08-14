package utils

import (
	"log"
	"os"
)

func GetEnv(key string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("The ENV '%s' must be defined!", key)
	}
	return val
}
