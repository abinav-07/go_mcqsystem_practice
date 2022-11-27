package models

import "time"

type Question struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Type      string    `json:"type"`
	TestID    uint      `json:"test_id"`
	Test      Test      `json:"test" gorm:"foreignKey:TestID"`
	Option    []Option  `json:"options" gorm:"foreignKey:QuestionID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`

	QuestionID string `gorm:"->" json:"question_id"` // Non-table field for alias
}

// Gives Table name of Model
func (u Question) TableName() string {
	return "questions"
}
