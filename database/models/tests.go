package models

import "time"

type Test struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsAvailable bool       `json:"is_available"`
	Question    []Question `json:"questions" gorm:"foreignKey:TestID"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   time.Time  `json:"-"`
}

// Gives Table name of Model
func (u Test) TableName() string {
	return "tests"
}
