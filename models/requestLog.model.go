package models

import (
	"time"

	"github.com/google/uuid"
)

type RequestLog struct {
	ID           int        `gorm:"primaryKey" json:"id"`
	RequestTime  time.Time  `gorm:"index" json:"requestTime"`
	ResponseTime time.Time  `gorm:"index" json:"responseTime"`
	UserID       *uuid.UUID `gorm:"type:uuid;index" json:"userId"`
	User         User       `gorm:"foreignKey:UserID;references:ID" json:"user"`
	Duration     float64    `gorm:"index" json:"duration"`
	Method       string     `gorm:"type:varchar(255);not null" json:"method"`
	Path         string     `gorm:"type:varchar(255);not null" json:"path"`
	Headers      string     `gorm:"type:text;not null" json:"headers"`
	Body         string     `gorm:"type:text;not null" json:"body"`
	Response     string     `gorm:"type:text;not null" json:"response"`
}
