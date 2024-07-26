package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func GetDBUrl() string {
	return os.Getenv("DATABASE_URL")
}

func GetCommitStartDate() string {
	return os.Getenv("COMMIT_START_DATE")
}

func GetCommitEndDate() string {
	return os.Getenv("COMMIT_END_DATE")
}
