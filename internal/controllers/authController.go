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
	authGroup.GET("/token-payload", services.AuthMiddleware, getTokenPayload)
	authGroup.POST("/refresh", refreshToken)
	authGroup.POST("/add-password", services.AuthMiddleware, addPassword)

	// GitHub OAuth
	authGroup.GET("/github/callback", services.GitHubCallback)
}

func addPassword(c *gin.Context) {
	var body struct {
		Password string `json:"password" bind:"required"`
	}
	var err error
	if err = c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password required",
		})
		return
	}
	var res *services.LoginResponce
	if id, exists := c.Get("userID"); !exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User id not found in token payload",
		})
		return
	} else {
		res, err = services.AddPassword(&services.AddPasswordDto{
			Password: body.Password,
			UserID:   id.(int64),
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, res)
	}
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

func getTokenPayload(c *gin.Context) {
	var userId, userName, userMail any
	var ok bool
	userId, ok = c.Get("userID")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user id not found in token"})
	}
	userName, ok = c.Get("userName")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user name not found in token"})
	}
	userMail, ok = c.Get("userMail")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user mail not found in token"})
	}
	c.JSON(http.StatusOK, gin.H{"id": userId, "name": userName, "mail": userMail})
}

func refreshToken(c *gin.Context) {
	var refreshBody struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&refreshBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var accessToken, refreshToken, err = services.RefreshTokens(refreshBody.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
	})
}
