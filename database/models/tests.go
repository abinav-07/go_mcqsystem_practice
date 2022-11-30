package models

import "time"

type Test struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title,omitempty"`
	Description string     `json:"description,omitempty"`
	IsAvailable bool       `json:"is_available,omitempty"`
	Question    []Question `json:"questions,omitempty" gorm:"foreignKey:TestID"`
	CreatedAt   time.Time  `json:"-,omitempty"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   time.Time  `json:"-"`
}

// Gives Table name of Model
func (u Test) TableName() string {
	return "tests"
}
