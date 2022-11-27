package services

import (
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/infrastructure"
)

type RoleService struct {
	repository infrastructure.Database
}

// Constructor
func NewRoleService(db infrastructure.Database) RoleService {
	return RoleService{
		repository: db,
	}
}

// Get  User By Id
func (c RoleService) GetById(ID uint) (models.Role, error) {
	role := models.Role{}

	return role, c.repository.DB.Where("id = ?", ID).First(&role).Error
}
