package initialize

import (
	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/internal/controllers"
)

func InitializeControllers(r *gin.Engine) {
	controllers.AddAuthControllers(r)
}
