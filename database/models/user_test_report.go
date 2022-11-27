package models

import "time"

type UserTestReport struct {
	ID        uint      `json:"id"`
	HasPassed bool      `json:"has_passed"`
	TestId    uint      `json:"test_id"`
	Test      Test      `json:"test" gorm:"foreignKey:TestId"`
	UserId    uint      `json:"user_id"`
	User      User      `json:"user" gorm:"foreignKey:UserId"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

// Gives Table name of Model
func (u UserTestReport) TableName() string {
	return "user_test_report"
}
