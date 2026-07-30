package app

import (
	"github.com/Delkira544/rakiduam/config"
	"github.com/Delkira544/rakiduam/internal/auth"
	ldapclient "github.com/Delkira544/rakiduam/internal/platform/ldap"
	"github.com/Delkira544/rakiduam/internal/projects"
	"github.com/Delkira544/rakiduam/internal/user"
	"github.com/jmoiron/sqlx"
)

// Handlers agrupa los handlers HTTP de cada feature ya ensamblados.
type Handlers struct {
	Auth    *auth.Handler
	Project *projects.Handler
}

// buildHandlers instancia repository -> service -> handler de cada feature,
// resolviendo a mano las dependencias cruzadas entre features (p.ej. auth
// necesita user.UserService). Ninguna feature importa este paquete.
func buildHandlers(sqlDB *sqlx.DB, ldapClient *ldapclient.Client, cfg config.Config) *Handlers {
	userRepo := user.NewUserRepository(sqlDB)
	userSvc := user.NewUserService(userRepo)

	authRepo := auth.NewTokenRepository(sqlDB)
	ldapRepo := auth.NewLDAPRepository(ldapClient, cfg.LDAP.BaseDN)
	authSvc := auth.NewService(ldapRepo, authRepo, &authUserAdapter{userSvc},
		cfg.App.JWTSecret, cfg.App.JWTExpiry, cfg.App.RefreshExpiry)

	projectsRepo := projects.NewRepository(sqlDB)
	projectsSvc := projects.NewService(projectsRepo)

	return &Handlers{
		Auth:    auth.NewHandler(authSvc),
		Project: projects.NewHandler(projectsSvc),
	}
}
