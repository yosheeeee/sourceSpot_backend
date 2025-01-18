package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID                int64
	Login             string
	Name              string
	Mail              string
	PasswordHash      string
	PasswordSalt      string
	IsGitHubConnected bool
	GitHubOAuthToken  string
	GitHubId          int64
}

type GitHubUser struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

type UserDto struct {
	ID   int64
	Name string
	Mail string
}

type UserRegisterDto struct {
	Login    string `json:"login" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLoginDto struct {
	Mail     string `json:"mail" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (user *User) ToDto() *UserDto {
	return &UserDto{
		ID:   user.ID,
		Name: user.Name,
		Mail: user.Mail,
	}
}
