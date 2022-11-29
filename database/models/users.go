package models

import "time"

type User struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	RoleID    uint      `json:"role_id"`
	Role      Role      `json:"role" gorm:"ForeignKey:RoleID"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	DeletedAt time.Time `json:"-"`

	FirebaseUID string `gorm:"->" json:"firebase_uid"` // Readonly FirebaseUID alias
}

// Gives Table name of Model
func (u User) TableName() string {
	return "users"
}
