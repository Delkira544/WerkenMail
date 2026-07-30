package auth

import (
	"context"
	"time"

	"github.com/Delkira544/rakiduam/internal/platform/logger"
	sharedauth "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Service interface {
	Login(c context.Context, req LoginRequest) (*LoginResponse, error)
}

type UserAdaptService interface {
	SyncFromLDAP(ctx context.Context, req *SyncUserRequest) (*SyncUserResponse, error)
}

type authService struct {
	ldapAuth      LDAPGateway
	tokenRepo     TokenRepository
	userSvc       UserAdaptService
	jwtSecret     string
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
}

func NewService(ldapAuth LDAPGateway, tokenRepo TokenRepository, userSvc UserAdaptService, jwtSecret, jwtExpiry, refExpiry string) Service {
	d, err := time.ParseDuration(jwtExpiry)
	if err != nil {
		d = time.Hour * 24 // default to 24 hours if parsing fails
	}

	refD, err := time.ParseDuration(refExpiry)
	if err != nil {
		refD = time.Hour * 24 * 7 // default to 7 days if parsing fails
	}

	return &authService{ldapAuth: ldapAuth, tokenRepo: tokenRepo, userSvc: userSvc,
		jwtSecret: jwtSecret, jwtExpiry: d, refreshExpiry: refD}
}

func (s *authService) Login(ctx context.Context, req LoginRequest) (*LoginResponse, error) {
	log := logger.FromContext(ctx)

	ldapUser, err := s.ldapAuth.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, errors.Internal("") // error de sistema
	}
	if ldapUser == nil {
		return nil, ErrInvalidCredentials // credenciales inválidas
	}

	localUser, err := s.userSvc.SyncFromLDAP(ctx, &SyncUserRequest{
		Username: ldapUser.Username,
		Name:     ldapUser.FullName,
		Email:    ldapUser.Email,
		Role:     ldapUser.Role.String(),
	})
	if err != nil {
		return nil, errors.Internal("failed to sync user")
	}

	now := time.Now()
	accessClaims := &sharedauth.Claims{
		Role: ldapUser.Role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   ldapUser.Username,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([]byte(s.jwtSecret))
	if err != nil {
		return nil, errors.Internal("failed to sign token")
	}

	rawToken, hash, err := generateRefreshToken()
	if err != nil {
		return nil, errors.Internal("failed to generate refresh token")
	}

	err = s.tokenRepo.Create(ctx, &RefreshToken{
		ID:        uuid.New(),
		UserID:    localUser.ID,
		TokenHash: hash,
		ExpiresAt: now.Add(s.refreshExpiry),
	})
	if err != nil {
		return nil, errors.Internal("failed to store refresh token")
	}

	log.Info("Authentication successful for user: ", zap.String("username", ldapUser.Username))
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawToken, // ← solo se ve 1 vez
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
		Role:         ldapUser.Role.String(),
	}, nil
}
