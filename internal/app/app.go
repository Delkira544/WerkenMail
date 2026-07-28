package app

import (
	"fmt"

	"github.com/Delkira544/rakiduam/config"
	"github.com/Delkira544/rakiduam/internal/platform/database"
	ldapclient "github.com/Delkira544/rakiduam/internal/platform/ldap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// App agrupa las dependencias de nivel raíz de la aplicación ya ensambladas.
type App struct {
	Cfg    config.Config
	Router *gin.Engine
	Pool   *pgxpool.Pool
}

// New es el composition root: crea la infraestructura (DB, LDAP), ensambla
// las features (repos -> services -> handlers) y arma el router. Devuelve
// una función cleanup para cerrar los recursos abiertos.
func New(cfg config.Config) (*App, func(), error) {
	pool, err := database.NewPool(cfg.Postgres)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	if err := database.RunMigrations(pool, "db/migrations"); err != nil {
		pool.Close()
		return nil, nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	sqlDB := database.NewSQLDB(pool)

	ldapClient, err := ldapclient.NewClient(cfg.LDAP)
	if err != nil {
		sqlDB.Close()
		pool.Close()
		return nil, nil, fmt.Errorf("failed to connect to LDAP: %w", err)
	}

	handlers := buildHandlers(sqlDB, ldapClient, cfg)
	router := NewRouter(handlers)

	cleanup := func() {
		sqlDB.Close()
		pool.Close()
		ldapClient.Close()
	}

	return &App{Cfg: cfg, Router: router, Pool: pool}, cleanup, nil
}
