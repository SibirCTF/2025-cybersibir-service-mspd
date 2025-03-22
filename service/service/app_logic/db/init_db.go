package db

import (
	"sibir2025/service/app_logic/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init_db() (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(sqlite.Open("MSPD2.db"), &gorm.Config{})
	if err != nil {
		return DB, err
	}
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		return DB, err
	}
	if err := DB.AutoMigrate(&models.Suspect{}); err != nil {
		return DB, err
	}
	if err := DB.AutoMigrate(&models.Claim{}); err != nil {
		return DB, err
	}
	return DB, nil
}
