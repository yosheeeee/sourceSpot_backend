package initialize

import (
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"github.com/yosheeeee/sourceSpot_baackend/database"
)

func InitializeApp(cfgPath string) error {
	if err := InitializeDBConnection(cfgPath); err != nil {
		return err
	}
	var r = gin.Default()
	r.Use(cors.Default())

	InitializeControllers(r)
	r.Run(fmt.Sprintf(":%d", config.AppConfig.Port))
	return nil
}

func InitializeDBConnection(cfgPath string) error {
	var err = config.LoadConfig(cfgPath)
	if err != nil {
		return err
	}
	if err = database.ConnectToDB(); err != nil {
		return err
	}
	return nil
}
