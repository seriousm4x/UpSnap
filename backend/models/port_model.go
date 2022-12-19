package models

import "gorm.io/gorm"

type Port struct {
	gorm.Model
	Number uint16 `gorm:"not null" json:"number" binding:"number"`
	Name   string `gorm:"not null" json:"name" binding:"name"`
}
