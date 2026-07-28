package user

type CreateUserRequest struct {
	LDAPUID  string `json:"ldap_uid" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
}
