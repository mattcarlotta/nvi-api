package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/mattcarlotta/nvi-api/utils"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Secret struct {
	ID           uuid.UUID     `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	UserId       uuid.UUID     `gorm:"type:uuid" json:"userId"`
	User         User          `gorm:"foreignKey:UserId;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Environments []Environment `gorm:"many2many:environment_secrets;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"environments"`
	Key          string        `gorm:"type:varchar(255);not null" json:"key"`
	Value        []byte        `gorm:"not null" json:"value"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

func (secret *Secret) BeforeCreate(tx *gorm.DB) (err error) {
	if encText, err := utils.CreateEncryptedText([]byte(secret.Value)); err == nil {
		tx.Statement.SetColumn("Value", encText)
	}
	return nil
}

func GetDupKeyinEnvs(secrets *[]Secret) string {
	var envNames string
	for _, secret := range *secrets {
		for _, env := range secret.Environments {
			if len(envNames) == 0 {
				envNames += env.Name
			} else {
				envNames += fmt.Sprintf(", %s", env.Name)
			}
		}
	}

	return envNames
}

func FindDupEnvNames(secrets *[]Secret, ids []uuid.UUID) string {
	var envNames string
	for _, secret := range *secrets {
		for _, env := range secret.Environments {
			for _, id := range ids {
				if env.ID == id {
					if len(envNames) == 0 {
						envNames += env.Name
					} else {
						envNames += fmt.Sprintf(", %s", env.Name)
					}
				}
			}
		}
	}

	return envNames
}

type SecretResult struct {
	ID           uuid.UUID      `json:"id"`
	UserId       uuid.UUID      `json:"userId"`
	Environments datatypes.JSON `json:"environments"`
	Key          string         `json:"key"`
	Value        []byte         `json:"value"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}
