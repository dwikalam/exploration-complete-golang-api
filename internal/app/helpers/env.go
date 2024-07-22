package helpers

import (
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	return nil
}

func GetEnvValue(envKey string, defaultValue string) string {
	v, ok := os.LookupEnv(envKey)

	if !ok {
		return defaultValue
	}

	return v
}
