package queries

import (
	"github.com/seriousm4x/UpSnap/database"
	"github.com/seriousm4x/UpSnap/models"
	"gorm.io/gorm"
)

func GetAllDevices(d *[]models.Device) error {
	if result := database.DB.Model(d).Find(d); result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateDevice(d *models.Device) error {
	if result := database.DB.Model(d).Create(&d); result.Error != nil {
		return result.Error
	}
	return nil
}

func GetOneDevice(d *models.Device, id int) error {
	result := database.DB.Model(d).Where("id = ?", id).Find(d)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func PatchDevice(changes *map[string]interface{}, id int) error {
	var device models.Device
	if err := GetOneDevice(&device, id); err != nil {
		return err
	}
	if result := database.DB.Model(device).Where("id = ?", id).Updates(changes); result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteDevice(d *models.Device, id int) error {
	var device models.Device
	if result := database.DB.Where("id = ?", id).Delete(&device); result.Error != nil {
		return result.Error
	}
	return nil
}
