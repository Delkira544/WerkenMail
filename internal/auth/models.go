package auth

type UserRole int

const (
	RoleStudent UserRole = iota
	RoleFunc
)

func (r UserRole) String() string {
	switch r {
	case RoleStudent:
		return "student"
	case RoleFunc:
		return "func"
	default:
		return "student"
	}
}

func ToRole(roleStr string) UserRole {
	switch roleStr {
	case "600":
		return RoleStudent
	case "500":
		return RoleFunc
	default:
		return RoleStudent // Default role if not recognized
	}
}

type User struct {
	Username string
	FullName string
	Email    string
	DN       string
	Role     UserRole
}
