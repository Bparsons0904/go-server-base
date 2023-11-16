package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`

	FirstName string `gorm:"type:varchar(255);not null" json:"firstName"`
	LastName  string `gorm:"type:varchar(255);not null" json:"lastName"`
	Email     string `gorm:"type:varchar(255);not null;unique" json:"email"`
	Username  string `gorm:"type:varchar(255);not null;unique" json:"username"`
	Password  string `gorm:"type:varchar(255);not null" json:"-"`
}
