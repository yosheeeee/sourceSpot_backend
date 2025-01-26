package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"github.com/yosheeeee/sourceSpot_baackend/database"
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
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error":   "Authorization code is expired or invalid. Please try again.",
				"details": string(oauth2Error.Body),
			})
			return
		}

		fmt.Println("Unexpected error:", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}
	// Получаем данные пользователя
	client := config.GetGitHubConfig().Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Логируем тело ответа
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	// Логируем в консоль сырой JSON
	fmt.Println("GitHub API Response:", string(bodyBytes))

	// Для повторного использования тела создаем новый io.Reader
	resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	// Декодируем JSON в структуру GitHubUser
	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Получаем email-адреса пользователя
	respEmails, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user emails"})
		return
	}
	defer respEmails.Body.Close()

	// Парсим список email-адресов
	var emails []struct {
		Email      string `json:"email"`
		Primary    bool   `json:"primary"`
		Verified   bool   `json:"verified"`
		Visibility string `json:"visibility"`
	}
	if err := json.NewDecoder(respEmails.Body).Decode(&emails); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode emails"})
		return
	}

	// Логируем email-адреса для проверки
	fmt.Println("GitHub User Emails:", emails)

	// Находим основной email
	var primaryEmail string
	for _, e := range emails {
		if e.Primary && e.Verified {
			primaryEmail = e.Email
			break
		}
	}

	// Если нашли email, добавляем его к данным пользователя
	if primaryEmail != "" {
		fmt.Println("Primary Email:", primaryEmail)
	} else {
		fmt.Println("No primary email found")
	}

	user.Email = primaryEmail
	var needPassword bool

	result, needPassword, err := getGitHubUserData(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	access_token, refresh_token, err := generateTokens(result)

	var respose = LoginResponce{
		User:         result,
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		NeedPassword: needPassword,
	}

	// Возвращаем данные пользователя
	c.JSON(http.StatusOK, respose)
}

func getGitHubUserData(gitHubUser *models.GitHubUser) (*models.UserDto, bool, error) {
	fmt.Println(gitHubUser.ID, gitHubUser.Name, gitHubUser.Email, gitHubUser.Login)
	var dbUser, err = FindUserByGitHubId(gitHubUser.ID)
	if err != nil {
		var userWithEmail, err = FindUserByMail(gitHubUser.Email)
		if err != nil {
			var existingUser, err = CreateUser(&models.UserRegisterDto{
				Login: gitHubUser.Login,
				Name:  gitHubUser.Name,
				Mail:  gitHubUser.Email,
			})
			if err != nil {
				return nil, false, err
			}
			existingUser.PasswordHash = ""
			existingUser.PasswordSalt = ""
			existingUser.AvatarPath = gitHubUser.AvatarURL
			existingUser.IsGitHubConnected = true
			existingUser.IsLocalAvatar = false
			existingUser.GitHubId = gitHubUser.ID
			existingUser.Mail = gitHubUser.Email
			database.DB.Save(&existingUser)
			return existingUser.ToDto(), true, nil
		} else {
			userWithEmail.AvatarPath = gitHubUser.AvatarURL
			userWithEmail.IsGitHubConnected = true
			userWithEmail.IsLocalAvatar = false
			userWithEmail.GitHubId = gitHubUser.ID
			userWithEmail.Mail = gitHubUser.Email
			database.DB.Save(&userWithEmail)
			return userWithEmail.ToDto(), false, nil
		}
	}
	fmt.Println("user found")
	return dbUser.ToDto(), dbUser.PasswordHash == "", nil
}
