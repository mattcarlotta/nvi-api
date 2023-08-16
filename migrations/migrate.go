package main

import (
	"fmt"
	"log"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

func main() {
	var db = database.GetConnection()
	var err = db.AutoMigrate(&models.User{})
	if err != nil {
		log.Fatalf("Unable to migrate User model: %s", err.Error())
	}
	err = db.AutoMigrate(&models.Environment{})
	if err != nil {
		log.Fatalf("Unable to migrate Environment model: %s", err.Error())
	}
	fmt.Println("🎉 Migration complete")
}
