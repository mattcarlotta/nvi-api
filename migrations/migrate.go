package main

import (
	"fmt"
	"log"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

func main() {
	db := database.GetConnection()
	if err := db.Migrator().DropTable(&models.User{}); err != nil {
		log.Fatalf("Unable to drop user table: %s", err.Error())
	}
	if err := db.Migrator().DropTable(&models.Environment{}); err != nil {
		log.Fatalf("Unable to environment table: %s", err.Error())
	}
	if err := db.Migrator().DropTable(&models.Secret{}); err != nil {
		log.Fatalf("Unable to secret table: %s", err.Error())
	}

	if err := db.AutoMigrate(&models.User{}, &models.Environment{}, &models.Secret{}); err != nil {
		log.Fatalf("Unable to migrate models: %s", err.Error())
	}

	// err := db.AutoMigrate(&models.User{})
	// if err != nil {
	// 	log.Fatalf("Unable to migrate User model: %s", err.Error())
	// }
	// err = db.AutoMigrate(&models.Environment{})
	// if err != nil {
	// 	log.Fatalf("Unable to migrate Environment model: %s", err.Error())
	// }
	// err = db.AutoMigrate(&models.Secret{})
	// if err != nil {
	// 	log.Fatalf("Unable to migrate Secret model: %s", err.Error())
	// }
	fmt.Println("🎉 Migration complete")
}
