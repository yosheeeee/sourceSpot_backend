package database

import (
	"fmt"

	"github.com/yosheeeee/sourceSpot_baackend/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() error {
	var err error
	var cfg = config.AppConfig
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d", cfg.Database.Host, cfg.Database.User, cfg.Database.Password, cfg.Database.Name, cfg.Database.Port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("Failed to connecto to database")
	}
	return nil
}
