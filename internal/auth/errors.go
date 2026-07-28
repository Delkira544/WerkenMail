package auth

import apperrors "github.com/Delkira544/rakiduam/internal/shared/errors"

var (
	ErrInvalidCredentials = apperrors.Unauthorized("invalid username or password")
	ErrSessionExpired     = apperrors.Unauthorized("session expired")
)
