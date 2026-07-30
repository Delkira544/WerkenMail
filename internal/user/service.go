package user

import (
	"context"

	authshared "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/google/uuid"
)

type UserService interface {
	Create(ctx context.Context, user *CreateUserRequest) error
	GetByUsername(ctx context.Context, ldapUID string) (*UserResponse, error)
	SyncFromLDAP(ctx context.Context, req *CreateUserRequest) (*UserResponse, error)
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
		Role:     authshared.ToRole(req.Role),
		Email:    req.Email,
	}

	err := s.repo.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) GetByUsername(ctx context.Context, ldapUID string) (*UserResponse, error) {
	user, err := s.repo.FindByUsername(ctx, ldapUID)
	if err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *userService) SyncFromLDAP(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	user := &User{
		ID:       uuid.New(),
		Username: req.Username,
		Name:     req.Name,
		Role:     authshared.ToRole(req.Role),
		Email:    req.Email,
	}
	u, err := s.repo.Upsert(ctx, user)
	if err != nil {
		return nil, err
	}
	return u.ToResponse(), nil
}
