package services

import (
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/infrastructure"
)

type UserTestReportService struct {
	repository infrastructure.Database
}

// Constructor
func NewUserTestService(repository infrastructure.Database) UserTestReportService {
	return UserTestReportService{
		repository: repository,
	}
}

// Create test report
func (ut UserTestReportService) Create(userReport models.UserTestReport) (models.UserTestReport, error) {
	return userReport, ut.repository.DB.Create(&userReport).Error
}

// Get user test report
func (ut UserTestReportService) GetTestReport(userID uint, testID uint) (*models.UserTestReport, error) {
	userTest := models.UserTestReport{}
	query := ut.repository.DB.Where("user_id = ? AND test_id = ?", userID, testID).First(&userTest)
	if query.Error != nil {
		return nil, query.Error
	}
	return &userTest, nil
}
