package main

import (
	"fmt"
	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

func main() {
	var db = database.GetDB()
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Environment{})
	fmt.Println("🎉 Migration complete")
}
