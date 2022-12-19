package models

import "gorm.io/gorm"

type Settings struct {
	gorm.Model
	SortBy               string `gorm:"not null; default:name" json:"sort_by" binding:"required"`
	ScanAddress          string `gorm:"not null" json:"scan_address"`
	PingInterval         uint   `gorm:"not null; default:5" json:"ping_interval"`
	NotificationsEnabled bool   `gorm:"not null; default:true" json:"notifications_enabled" binding:"required"`
}
