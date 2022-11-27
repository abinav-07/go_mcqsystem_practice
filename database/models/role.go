package models

import "time"

type Role struct {
	ID        uint      `json:"id"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`
}

// Gives Table name of Model
func (r Role) TableName() string {
	return "roles"
}
