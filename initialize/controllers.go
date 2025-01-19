package initialize

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/internal/controllers"
)

func InitializeControllers(r *gin.Engine) {
	disableCORS(r)
	controllers.AddAuthControllers(r)
	controllers.AddUserControllers(r)
}

func disableCORS(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
}
