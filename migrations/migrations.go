package main

import (
	"github.com/yosheeeee/sourceSpot_baackend/database"
	"github.com/yosheeeee/sourceSpot_baackend/initialize"
	"github.com/yosheeeee/sourceSpot_baackend/internal/models"
)

func init() {
	initialize.InitializeDBConnection("../config.yaml")
}

func main() {
	database.DB.AutoMigrate(&models.User{})
}
