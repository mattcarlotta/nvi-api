package database

import (
	"fmt"
	"log"
	"os"

	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func CreateConnection() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		utils.GetEnv("DB_HOST"),
		utils.GetEnv("DB_USER"),
		utils.GetEnv("DB_PASSWORD"),
		utils.GetEnv("DB_NAME"),
		utils.GetEnv("DB_PORT"),
	)

	logType := logger.Info
	if os.Getenv("IN_PRODUCTION") == "true" || os.Getenv("IN_TESTING") == "true" {
		logType = logger.Silent
	}

	dbc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logType),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	fmt.Println("🔌 Successfully connected to the database")
	DB = dbc
	return DB
}

func GetConnection() *gorm.DB {
	if DB == nil {
		return CreateConnection()
	}
	return DB
}
