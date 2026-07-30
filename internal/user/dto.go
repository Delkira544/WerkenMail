package user

type CreateUserRequest struct {
	Username string `json:"ldap_uid" validate:"required"`
	Name     string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Role     string `json:"role" validate:"required"`
}

type UserResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at,omitempty"`
}
