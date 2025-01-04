package services

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yosheeeee/sourceSpot_baackend/config"
	"github.com/yosheeeee/sourceSpot_baackend/internal/models"
)

type LoginResponce struct {
	User         models.UserDto
	AccessToken  string
	RefreshToken string
}

func RegisterUser(createDto *models.UserRegisterDto) (*LoginResponce, error) {
	var user *models.UserDto
	var err error
	if user, err = CreateUser(createDto); err != nil {
		return nil, err
	}

	var accessToken, refreshToken string
	accessToken, refreshToken, err = generateTokens(user)
	if err != nil {
		return nil, err
	}

	return &LoginResponce{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
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
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func generateTokens(user *models.UserDto) (accessToken string, refreshToken string, err error) {
	accessClaims := &models.TokenClaims{
		ID: int(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}
	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(config.AppConfig.JWTSecretKey)
	if err != nil {
		return "", "", err
	}

	// Refresh Token
	refreshClaims := &models.TokenClaims{
		ID: int(user.ID),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
		},
	}
	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(config.AppConfig.JWTSecretKey)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// func getTokenPayload(token string) (*models.TokenClaims, error) {

// }
