package models

import (
	"time"

	"github.com/google/uuid"
)

type Project struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primary_key" json:"id"`
	Name      string    `gorm:"type:varchar(255);index:project_index;not null" json:"name"`
	UserID    uuid.UUID `gorm:"type:uuid;index:project_index" json:"userID"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ReqProject struct {
	Name string `json:"name" validate:"required,name,lte=255"`
}

type ReqUpdateProject struct {
	ID          string `json:"id" validate:"required,uuid"`
	UpdatedName string `json:"name" validate:"required,name,lte=255"`
}
