package dtos

type UserClaimMetaData struct {
	UserRole string
	UserID   uint
	IsAdmin  bool
}

type AdminUserUpdateData struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	RoleID    uint   `json:"role_id,omitempty"`
}
