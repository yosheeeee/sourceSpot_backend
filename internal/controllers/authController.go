package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/internal/models"
	"github.com/yosheeeee/sourceSpot_baackend/internal/services"
)

func AddAuthControllers(r *gin.Engine) {
	var authGroup = r.Group("/auth")
	authGroup.POST("/login", loginUser)
	authGroup.POST("/register", registerUser)
}

func loginUser(c *gin.Context) {
	var req *models.UserLoginDto
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var res *services.LoginResponce
	if res, err = services.LoginUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &res)
}

func registerUser(c *gin.Context) {
	var req *models.UserRegisterDto
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	var res *services.LoginResponce
	if res, err = services.RegisterUser(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &res)
}
