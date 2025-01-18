package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"github.com/yosheeeee/sourceSpot_baackend/internal/models"
	"golang.org/x/oauth2"
)

func GitHubCallback(c *gin.Context) {
	// Получаем код авторизации из callback-а
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not found"})
		return
	}

	token, err := config.GetGitHubConfig().Exchange(context.Background(), code)
	if err != nil {
		if oauth2Error, ok := err.(*oauth2.RetrieveError); ok {
			fmt.Printf("Error details: %s\n", string(oauth2Error.Body))
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Authorization code is expired or invalid. Please try again.",
				"details": string(oauth2Error.Body),
			})
			return
		}

		fmt.Println("Unexpected error:", err)
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

	// Декодируем JSON в структуру GitHubUser
	var user models.GitHubUser
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	result, err := getGitHubUserData(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	access_token, refresh_token, err := generateTokens(result)

	var respose = LoginResponce{
		User:         result,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	}

	// Возвращаем данные пользователя
	c.JSON(http.StatusOK, respose)
}

func getGitHubUserData(gitHubUser *models.GitHubUser) (*models.UserDto, error) {
	var dbUser, err = FindUserByGitHubId(gitHubUser.ID)
	if err != nil {
		var existingUser, err = CreateUser(&models.UserRegisterDto{
			Login: gitHubUser.Login,
			Name:  gitHubUser.Name,
			Mail:  gitHubUser.Email,
		})
		if err != nil {
			return nil, err
		}
		return existingUser, nil
	}
	return dbUser.ToDto(), nil
}
