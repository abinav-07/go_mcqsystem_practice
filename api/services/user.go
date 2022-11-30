package services

import (
	"errors"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/dtos"
	"github/abinav-07/mcq-test/infrastructure"

	"gorm.io/gorm"
)

type UserService struct {
	repository      infrastructure.Database
	firebaseService FirebaseService
}

// Constructor
func NewUserService(database infrastructure.Database, firebaseService FirebaseService) UserService {
	return UserService{
		repository:      database,
		firebaseService: firebaseService,
	}
}

// WithTrx -> enables repository with transaction
func (c UserService) WithTrx(trxHandle *gorm.DB) (UserService, error) {
	if trxHandle == nil {
		return c, errors.New("Transaction DB not found")
	}

	c.repository.DB = trxHandle
	return c, nil
}

// Create User
func (c UserService) CreateUserWithFB(user models.User, claimData dtos.UserClaimMetaData) (*models.User, error) {
	createdUser, createUserErr := c.Create(user)

	if createUserErr != nil {
		return nil, createUserErr
	}

	claimData.UserID = createdUser.ID

	//Create Firebase user
	fb_uid, err := c.firebaseService.GetCreateOrUpdateFirebaseUser(createdUser, claimData)
	createdUser.FirebaseUID = fb_uid

	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

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

// Create Test
func (t UserService) DeleteById(ID uint) (*models.User, error) {
	user := models.User{}
	err := t.repository.DB.Model(&user).Where("id = ?", ID).Find(&user).Delete(&user).Error

	return &user, err
}
