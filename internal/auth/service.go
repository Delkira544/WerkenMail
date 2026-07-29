package auth

import (
	"context"
	"time"

	sharedauth "github.com/Delkira544/rakiduam/internal/shared/auth"
	"github.com/Delkira544/rakiduam/internal/shared/errors"
	"github.com/Delkira544/rakiduam/internal/user"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Service interface {
	Login(c context.Context, req LoginRequest) (*LoginResponse, error)
}

type authService struct {
	ldapAuth      LDAPRepository
	tokenRepo     TokenRepository
	userSvc       user.UserService
	jwtSecret     string
	jwtExpiry     time.Duration
	refreshExpiry time.Duration
}

func NewService(ldapAuth LDAPRepository, tokenRepo TokenRepository, userSvc user.UserService, jwtSecret, jwtExpiry, refExpiry string) Service {
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
	// ── 1. Autenticar contra LDAP ──
	ldapUser, err := s.ldapAuth.Authenticate(req.Username, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials // error de sistema
	}
	if ldapUser == nil {
		return nil, ErrInvalidCredentials // credenciales inválidas
	}
	// ── 2. Upsert en users local (cachear role) ──
	_ = s.userSvc.SyncFromLDAP(ctx, &user.CreateUserRequest{
		Username: ldapUser.Username,
		Name:     ldapUser.FullName,
		Email:    ldapUser.Email,
		Role:     ldapUser.Role.String(),
	})
	// ── 3. Generar access_token ──
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
	// ── 4. Generar refresh_token aleatorio ──
	rawToken, hash, err := generateRefreshToken()
	if err != nil {
		return nil, errors.Internal("failed to generate refresh token")
	}
	// ── 5. Persistir hash en DB ──
	err = s.tokenRepo.Create(ctx, &RefreshToken{
		ID:        uuid.New(),
		UserID:    ldapUser.Username,
		TokenHash: hash,
		ExpiresAt: now.Add(s.refreshExpiry),
	})
	if err != nil {
		return nil, errors.Internal("failed to store refresh token")
	}
	// ── 6. Responder ──
	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawToken, // ← solo se ve 1 vez
		ExpiresIn:    int64(s.jwtExpiry.Seconds()),
	}, nil
}
