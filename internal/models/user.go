package models

import (
	"time"
)

type User struct {
	ID            uint      `gorm:"primaryKey"`
	Email         string    `gorm:"unique;not null"`
	Username      string    `gorm:"unique;not null;size:20"`
	DisplayName   string    `gorm:"size:50"`
	Password      string    `gorm:"not null"`
	IsAdmin       bool      `gorm:"default:false"`
	CreatedAt     time.Time
	SetupStatus   string    `gorm:"default:'not_started'"`
	Language      string    `gorm:"default:'en'"`
	StatusPrivacy string    `gorm:"default:'friends'"`
	
	CurrentStatus string     `gorm:"default:'offline'"`
	CurrentGame   string     `gorm:"default:''"`
	ClientType    string     `gorm:"default:'web'"`
	LastSeenAt    *time.Time `gorm:"default:null"`

	Pal Pal `gorm:"embedded;embeddedPrefix:pal_"`
}

type Pal struct {
	Config    string    `gorm:"type:json"`
}