package main

import (
	"fmt"
	"log"

	"github.com/mattcarlotta/nvi-api/database"
	"github.com/mattcarlotta/nvi-api/models"
)

func main() {
	db := database.GetConnection()
	db.Migrator().DropTable(&models.User{})
	db.Migrator().DropTable(&models.Environment{})
	db.Migrator().DropTable(&models.Secret{})

	err := db.AutoMigrate(&models.User{}, &models.Environment{}, &models.Secret{})
	if err != nil {
		log.Fatalf("Unable to migrate User model: %s", err.Error())
	}
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
