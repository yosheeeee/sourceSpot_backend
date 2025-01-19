package services

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"github.com/yosheeeee/sourceSpot_baackend/internal/models"
)

type LoginResponce struct {
	User         *models.UserDto
	AccessToken  string
	RefreshToken string
}

func RegisterUser(createDto *models.UserRegisterDto) (*LoginResponce, error) {
	var user *models.User
	var err error
	if user, err = CreateUser(createDto); err != nil {
		return nil, err
	}

	var accessToken, refreshToken string
	accessToken, refreshToken, err = generateTokens(&models.UserDto{
		ID:   user.ID,
		Name: user.Name,
		Mail: user.Mail,
	})
	if err != nil {
		return nil, err
	}

	return &LoginResponce{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToDto(),
	}, nil
}

func LoginUser(loginDto *models.UserLoginDto) (*LoginResponce, error) {
	var user, err = FindUserByMail(loginDto.Mail)
	if err != nil {
		return nil, err
	}
	var hashResult bool
	hashResult, err = VerifyUserPassword(loginDto.Password, user.PasswordHash, user.PasswordSalt)
	if err != nil {
		return nil, err
	}
	if !hashResult {
		return nil, fmt.Errorf("Invalid password")
	}
	var accessToken, refreshToken string
	if accessToken, refreshToken, err = generateTokens(user.ToDto()); err != nil {
		return nil, err
	}
	return &LoginResponce{
		User:         user.ToDto(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateTokens(user *models.UserDto) (accessToken string, refreshToken string, err error) {
	accessClaims := &models.TokenClaims{
		ID:   user.ID,
		Name: user.Name,
		Mail: user.Mail,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(config.GetJWTSecretKey())
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := &models.TokenClaims{
		ID:   user.ID,
		Mail: user.Mail,
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(config.GetJWTSecretKey())
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func AuthMiddleware(c *gin.Context) {
	// Получаем токен из заголовка Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	// Проверяем формат заголовка (должен начинаться с "Bearer ")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
		return
	}

	// Парсим токен
	claims, err := getTokenPayload(tokenString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	// Добавляем информацию о пользователе в контекст
	c.Set("userID", claims.ID)
	c.Set("userMail", claims.Mail)
	c.Set("userName", claims.Name)

	// Передаем управление следующему обработчику
	c.Next()
}

func getTokenPayload(token string) (*models.TokenClaims, error) {
	parsedToken, err := jwt.ParseWithClaims(token, &models.TokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Проверяем алгоритм подписи
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return config.GetJWTSecretKey(), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	claims, ok := parsedToken.Claims.(*models.TokenClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Проверяем истечение срока действия токена
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}

func RefreshTokens(refreshToken string) (accessToken string, newRefreshToken string, err error) {
	// Парсим refresh-токен
	claims, err := getTokenPayload(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid or expired refresh token: %v", err)
	}

	// Проверяем истечение срока действия refresh-токена
	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return "", "", fmt.Errorf("refresh token has expired")
	}

	// Генерируем новые токены
	user := &models.UserDto{
		ID:   int64(claims.ID),
		Name: claims.Name,
		Mail: claims.Mail,
	}
	accessToken, newRefreshToken, err = generateTokens(user)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate new tokens: %v", err)
	}
	return accessToken, newRefreshToken, nil
}
