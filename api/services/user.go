package services

import (
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/infrastructure"
)

type UserService struct {
	repository infrastructure.Database
}

// Constructor
func NewUserService(database infrastructure.Database) UserService {
	return UserService{
		repository: database,
	}
}

// Create User
func (c UserService) Create(user models.User) (*models.User, error) {
	return &user, c.repository.DB.Create(&user).Preload("Role").Error
}

// Get  User By Id
func (c UserService) GetById(ID uint) (*models.User, error) {
	user := models.User{}
	return &user, c.repository.DB.Model(&models.User{}).Where("id = ?", ID).Preload("Role").First(&user).Error
}

// Get  User By Email
func (c UserService) GetByEmail(email string) (*models.User, error) {
	user := models.User{}

	return &user, c.repository.DB.Where("email = ?", email).Preload("Role").First(&user).Error
}
