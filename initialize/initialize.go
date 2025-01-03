package initialize

import (
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"github.com/yosheeeee/sourceSpot_baackend/database"
)

func InitializeApp(cfgPath string) error {
	if err := InitializeDBConnection(cfgPath); err != nil {
		return err
	}
	return nil
}

func InitializeDBConnection(cfgPath string) error {
	var cfg, err = config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	if err = database.ConnectToDB(cfg); err != nil {
		return err
	}
	return nil
}
