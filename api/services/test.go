package services

import (
	"errors"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/infrastructure"

	"gorm.io/gorm"
)

type TestService struct {
	repository infrastructure.Database
}

// contructor
func NewTestService(database infrastructure.Database) TestService {
	return TestService{
		repository: database,
	}
}

// WithTrx -> enables repository with transaction
func (c TestService) WithTrx(trxHandle *gorm.DB) (TestService, error) {
	if trxHandle == nil {
		return c, errors.New("Transaction DB not found")
	}

	c.repository.DB = trxHandle
	return c, nil
}

// Create Test
func (t TestService) Create(test models.Test) (*models.Test, error) {
	return &test, t.repository.DB.Create(&test).Error
}

// Get  Test By Id
func (t TestService) GetById(ID uint) (*models.Test, error) {
	test := models.Test{}
	return &test, t.repository.DB.Preload("Question").Where("id = ?", ID).First(&test).Error
}

// Get  Test By Id
func (t TestService) GetByName(name string) (*models.Test, error) {
	test := models.Test{}
	return &test, t.repository.DB.Where("title LIKE ?", name).First(&test).Error
}

func (t TestService) GetTestByQuery(queryParams models.Test) ([]models.Test, error) {
	test := []models.Test{}
	queryBuilder := t.repository.DB
	queryBuilder = queryBuilder.Model(&models.Test{})

	if queryParams.IsAvailable == true {
		queryBuilder.Where("is_available = ?", true)
	}

	return test, queryBuilder.Find(&test).Error

}

func (t TestService) UpdateOneTest(testID uint, updateTest map[string]interface{}) (models.Test, error) {
	test := models.Test{}
	return test, t.repository.DB.Model(&models.Test{}).Where("id = ?", testID).Updates(updateTest).Find(&test).Error
}
