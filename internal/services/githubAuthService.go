package services

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"golang.org/x/oauth2"
)

func GitHubLogin(c *gin.Context) {
	// Генерируем URL для авторизации
	url := config.GetGitHubConfig().AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GitHubCallback(c *gin.Context) {
	// Получаем код авторизации из callback-а
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	// Обмениваем код на токен
	token, err := config.GetGitHubConfig().Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Получаем данные пользователя
	client := config.GetGitHubConfig().Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Возвращаем данные пользователя
	c.JSON(http.StatusOK, gin.H{"message": "GitHub login successful"})
}
