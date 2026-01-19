package models

import (
	"time"
)

type User struct {
	ID          uint      `gorm:"primaryKey"`
	Username    string    `gorm:"unique;not null;size:20"`
	DisplayName string    `gorm:"size:50"`
	Password    string    `gorm:"not null"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Status      string    `gorm:"default:'offline'"`
	SetupStatus string    `gorm:"default:'not_started'"`
	Language    string    `gorm:"default:'en'"`

	Pal       Pal       `gorm:"embedded;embeddedPrefix:pal_"`
}

type Pal struct {
	Config    string    `gorm:"type:json"`
}