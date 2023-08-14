package database

import (
	"fmt"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		utils.GetEnv("DB_HOST"),
		utils.GetEnv("DB_USER"),
		utils.GetEnv("DB_PASSWORD"),
		utils.GetEnv("DB_NAME"),
		utils.GetEnv("DB_PORT"),
	)

	dbc, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to the database")
	}

	fmt.Println("🔌 Successfully connected to the database")
	DB = dbc
	return DB
}

func GetDB() *gorm.DB {
	if DB == nil {
		return ConnectDB()
	} else {
		return DB
	}
}
