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
	UserID       uuid.UUID     `gorm:"type:uuid" json:"userID"`
	User         User          `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Environments []Environment `gorm:"many2many:environment_secrets;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"environments"`
	Key          string        `gorm:"type:varchar(255);not null" json:"key"`
	Value        []byte        `gorm:"not null" json:"value"`
	Nonce        []byte        `gorm:"not null" json:"nonce"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

func (secret *Secret) BeforeCreate(tx *gorm.DB) (err error) {
	if encText, nonce, err := utils.CreateEncryptedSecretValue(secret.Value); err == nil {
		tx.Statement.SetColumn("Value", encText)
		tx.Statement.SetColumn("Nonce", nonce)
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
	UserID       uuid.UUID      `json:"userID"`
	Environments datatypes.JSON `json:"environments"`
	Key          string         `json:"key"`
	Value        []byte         `json:"value"`
	Nonce        []byte         `json:"-"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

type ReqCreateSecret struct {
	ProjectID      string   `json:"projectID" validate:"required,uuid"`
	EnvironmentIDs []string `json:"environmentIDs" validate:"uuidarray"`
	Key            string   `json:"key" validate:"required,gte=2,lte=255"`
	Value          string   `json:"value" validate:"required,lte=5000"`
}

type ReqUpdateSecret struct {
	ID             string   `json:"id" validate:"required,uuid"`
	EnvironmentIDs []string `json:"environmentIDs" validate:"uuidarray"`
	Key            string   `json:"key" validate:"required,gte=2,lte=255"`
	Value          string   `json:"value" validate:"required,lte=5000"`
}
