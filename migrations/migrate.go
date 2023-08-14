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
	var DB = database.GetDB()
	DB.AutoMigrate(&models.User{})
	fmt.Println("? Migration complete")
}
