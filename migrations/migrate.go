package main

import (
	"fmt"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

func init() {
	database.ConnectDB()
}

func main() {
	database.DB.AutoMigrate(&models.User{})
	fmt.Println("? Migration complete")
}
