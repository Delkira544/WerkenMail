package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, user *CreateUserRequest) error
	GetUserByLDAPUID(ctx context.Context, ldapUID string) (*User, error)
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) CreateUser(ctx context.Context, req *CreateUserRequest) error {
	user := &User{
		ID:       uuid.New(),
		LDAPUID:  req.LDAPUID,
		FullName: req.FullName,
		Email:    req.Email,
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetUserByLDAPUID(ctx context.Context, ldapUID string) (*User, error) {
	user, err := s.repo.FindByLDAPUID(ctx, ldapUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
