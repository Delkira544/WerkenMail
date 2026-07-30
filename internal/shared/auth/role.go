package auth

type Role string

const (
	RoleStudent Role = "student"
	RoleFunc    Role = "func"
)

func (r Role) String() string {
	return string(r)
}

func (r Role) IsStudent() bool {
	return r == RoleStudent
}

func (r Role) IsFunc() bool {
	return r == RoleFunc
}

func ToRole(role string) Role {
	switch role {
	case "student":
		return RoleStudent
	case "func":
		return RoleFunc
	default:
		return RoleStudent // Default role if not recognized
	}
}
