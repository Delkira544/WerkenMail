package user

type CreateUserRequest struct {
	Username string `json:"ldap_uid" validate:"required"`
	Name     string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}
