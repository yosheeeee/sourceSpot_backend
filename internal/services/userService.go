package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/yosheeeee/sourceSpot_baackend/database"
	"github.com/yosheeeee/sourceSpot_baackend/internal/models"
	"gorm.io/gorm"
)

// функция создания пользователя
func CreateUser(createDto *models.UserRegisterDto) (*models.UserDto, error) {
	var existingUser models.User
	database.DB.Where(models.User{Name: createDto.Name}).Or(models.User{Mail: createDto.Mail}).Find(&existingUser)
	if existingUser.ID != 0 {
		return nil, fmt.Errorf("User with this name or email is already exists")
	}
	var hash, salt string
	var err error
	if hash, salt, err = generatePasswordHash(createDto.Password); err != nil {
		return nil, err
	}
	var user = models.User{
		Name:         createDto.Name,
		Mail:         createDto.Mail,
		PasswordHash: hash,
		PasswordSalt: salt,
	}
	database.DB.Create(&user)
	return user.ToDto(), nil
}

// сравнение паролей
func VerifyUserPassword(password string, hash string, salt string) (bool, error) {
	var passwdHash, err = hashPassword(password, salt)
	if err != nil {
		return false, err
	}
	return passwdHash == hash, nil
}

// создание хэша по паролю и соли
func hashPassword(password string, salt string) (string, error) {
	var data = password + salt
	hash := sha256.Sum256([]byte(data))
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// генерация хэша и соли по паролю
func generatePasswordHash(password string) (string, string, error) {
	var salt, err = generateSalt()
	if err != nil {
		return "", "", err
	}
	var hash string
	if hash, err = hashPassword(password, salt); err != nil {
		return "", "", err
	}
	return hash, salt, nil
}

// генерирование соли
func generateSalt() (string, error) {
	salt := make([]byte, 10)
	_, err := rand.Read(salt)
	if err != nil {
		return "", fmt.Errorf("Error while generating hash salt: %v", err)
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

// функция поиска пользователя по id
func FindUserById(id int64) (*models.User, error) {
	var user models.User
	// Ищем пользователя по ID
	if err := database.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %d not found", id)
		}
		return nil, err // Возвращаем другие ошибки базы данных
	}
	return &user, nil
}

// поиск пользователя по email
func FindUserByMail(mail string) (*models.User, error) {
	var user models.User
	// Ищем пользователя по ID
	if err := database.DB.First(&user, models.User{Mail: mail}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with mail %s not found", mail)
		}
		return nil, err // Возвращаем другие ошибки базы данных
	}
	return &user, nil
}
