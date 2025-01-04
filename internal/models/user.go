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
	PaswordSalt  string
}

type UserDto struct {
	ID   int64
	Name string
	Mail string
}

type UserCreateDto struct {
	Name     string
	Mail     string
	Password string
}

func (user *User) ToDto() *UserDto {
	return &UserDto{
		ID:   user.ID,
		Name: user.Name,
		Mail: user.Mail,
	}
}
