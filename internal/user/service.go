package user

import (
	"context"

	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user *CreateUserRequest) error
	GetByUsername(ctx context.Context, ldapUID string) (*User, error)
	SyncFromLDAP(ctx context.Context, req *CreateUserRequest) error
}

type userService struct {
	repo UserRepository
}

func NewUserService(repo UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, req *CreateUserRequest) error {
	user := &User{
		ID:       uuid.New(),
		Username: req.Username,
		Name:     req.Name,
		Role:     req.Role,
		Email:    req.Email,
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetByUsername(ctx context.Context, ldapUID string) (*User, error) {
	user, err := s.repo.FindByUsername(ctx, ldapUID)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (s *userService) SyncFromLDAP(ctx context.Context, req *CreateUserRequest) error {
	user := &User{
		ID:       uuid.New(),
		Username: req.Username,
		Name:     req.Name,
		Role:     req.Role,
		Email:    req.Email,
	}
	return s.repo.Upsert(ctx, user)
}
