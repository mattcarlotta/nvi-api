package models

import (
	"time"

	"github.com/google/uuid"
)

type Environment struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	ProjectID uuid.UUID `gorm:"type:uuid;index:env_index" json:"projectID"`
	Project   Project   `gorm:"foreignKey:ProjectID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	UserID    uuid.UUID `gorm:"type:uuid;index:env_index" json:"userID"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	Name      string    `gorm:"type:varchar(255);index:env_index;not null" json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ReqCreateEnv struct {
	Name      string `json:"name" validate:"required,name,lte=255"`
	ProjectID string `json:"projectID" validate:"uuid"`
}

type ReqUpdateEnv struct {
	ID          string `json:"id" validate:"required,uuid"`
	ProjectID   string `json:"projectID" validate:"uuid"`
	UpdatedName string `json:"updatedName" validate:"required,name,lte=255"`
}
