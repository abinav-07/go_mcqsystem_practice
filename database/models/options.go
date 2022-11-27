package models

import "time"

type Option struct {
	ID         uint      `json:"id"`
	Title      string    `json:"title"`
	IsCorrect  bool      `json:"is_correct"`
	QuestionID uint      `json:"question_id"`
	Question   Question  `json:"-" gorm:"foreignKey:QuestionID"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
	DeletedAt  time.Time `json:"-"`

	OptionId string `gorm:"->" json:"option_id"` // Non-table field for alias
}

// Gives Table name of Model
func (u Option) TableName() string {
	return "options"
}
