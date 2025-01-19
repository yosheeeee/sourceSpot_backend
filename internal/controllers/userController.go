package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/internal/services"
)

func AddUserControllers(r *gin.Engine) {
	var userGroup = r.Group("/user")
	userGroup.GET("/header-data", services.AuthMiddleware, GetUserData)
}

func GetUserData(c *gin.Context) {
	var userId, exists = c.Get("userID")
	if !exists {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "user id not found in claims",
		})
		return
	}
	var userIdint, ok = userId.(int64)
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "user id is not int",
		})
		return
	}
	var user, err = services.FindUserById(userIdint)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	c.JSON(http.StatusOK, user.ToDto())
}
