package models

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name        string `gorm:"not null; size:100" json:"name" binding:"required"`
	IP          string `gorm:"not null" json:"ip" binding:"required"`
	Mac         string `gorm:"not null; size:17" json:"mac" binding:"required"`
	Netmask     string `gorm:"not null; size:15; default:255.255.255.0" json:"netmask" binding:"required"`
	Link        string `json:"link"`
	Ports       []Port `gorm:"many2many:device_ports;" json:"ports"`
	ShutdownCmd string `json:"shutdown_cmd"`
}
