package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID           int64
	Name         string
	Mail         string
	PasswordHash string
	PasswordSalt string
}

type UserDto struct {
	ID   int64
	Name string
	Mail string
}

type UserRegisterDto struct {
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
