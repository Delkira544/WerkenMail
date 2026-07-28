package auth

import (
	"time"

	sharedauth "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/Delkira544/rakiduam/internal/user"
	"github.com/golang-jwt/jwt/v5"
)

type Service interface {
	Login(req LoginRequest) (*Session, error)
}

type authService struct {
	authRepo  Repository
	userRepo  user.UserService
	jwtSecret string
	jwtExpiry time.Duration
}

func NewService(repo Repository, userRepo user.UserService, jwtSecret, jwtExpiry string) Service {
	d, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		d = time.Hour * 24 // default to 24 hours if parsing fails
	}
	return &authService{authRepo: repo, userRepo: userRepo, jwtSecret: jwtSecret, jwtExpiry: d}
}

func (s *authService) Login(req LoginRequest) (*Session, error) {
	ldapUser, err := s.authRepo.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	now := time.Now()
	expiresAt := now.Add(s.jwtExpiry)

	claims := &sharedauth.Claims{
		Role: ldapUser.Role.String(), // This should be fetched
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   ldapUser.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, errors.Internal("failed to sign token")
	}

	return &Session{
		Token:     signedToken,
		ExpiresAt: expiresAt.Unix(),
		UserID:    ldapUser.Username,
		Role:      ldapUser.Role.String(),
	}, nil
}
