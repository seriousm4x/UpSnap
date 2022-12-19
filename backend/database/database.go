package database

import (
	"fmt"

	"github.com/seriousm4x/UpSnap/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&models.Settings{}, &models.Port{}, &models.Device{})

	DB = db
}

func Close() error {
	db, err := DB.DB()
	if err != nil {
		return fmt.Errorf("gorm.DB get database: %v", err)
	}
	return db.Close()
}
