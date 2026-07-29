package config

import (
	"strings"

	"github.com/spf13/viper"
)

func LoadConfig() (Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath("config")
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return Config{}, err
	}

	bindEnvVars(v)

	return *buildConfig(v), nil
}

func bindEnvVars(v *viper.Viper) {

	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	bindings := map[string]string{
		// App
		"app.name":                 "APP_NAME",
		"app.env":                  "APP_ENV",
		"app.port":                 "APP_PORT",
		"app.debug":                "APP_DEBUG",
		"app.timezone":             "TIMEZONE",
		"app.jwtSecret":            "JWT_SECRET",
		"app.jwtExpiration":        "JWT_EXPIRATION",
		"app.jwtRefreshExpiration": "JWT_REFRESH_EXPIRATION",

		// LDAP
		"ldap.host":     "LDAP_HOST",
		"ldap.port":     "LDAP_PORT",
		"ldap.bindDN":   "LDAP_BIND_DN",
		"ldap.baseDN":   "LDAP_BASE_DN",
		"ldap.bindPass": "LDAP_BIND_PASS",
		"ldap.useTLS":   "LDAP_USE_TLS",

		// PostgreSQL
		"postgres.host":         "DATABASE_HOST",
		"postgres.port":         "DATABASE_PORT",
		"postgres.database":     "DATABASE_NAME",
		"postgres.username":     "DATABASE_USERNAME",
		"postgres.password":     "DATABASE_PASSWORD",
		"postgres.sslMode":      "DATABASE_SSL_MODE",
		"postgres.maxOpenConns": "DATABASE_MAX_OPEN_CONNS",

		// Logger
		"logger.level":  "LOGGER_LEVEL",
		"logger.format": "LOGGER_FORMAT",
		"logger.output": "LOGGER_OUTPUT",
	}

	for key, env := range bindings {
		_ = v.BindEnv(key, env)
	}
}

func buildConfig(v *viper.Viper) *Config {
	return &Config{
		App: AppConfig{
			Name:     v.GetString("app.name"),
			Env:      v.GetString("app.env"),
			Port:     v.GetInt("app.port"),
			Debug:    v.GetBool("app.debug"),
			Timezone: v.GetString("app.timezone"),
		},
		LDAP: LDAPConfig{
			Host:     v.GetString("ldap.host"),
			Port:     v.GetInt("ldap.port"),
			BaseDN:   v.GetString("ldap.baseDN"),
			BindDN:   v.GetString("ldap.bindDN"),
			BindPass: v.GetString("ldap.bindPass"),
			UseTLS:   v.GetBool("ldap.useTLS"),
		},
		Postgres: PostgresConfig{
			Host:         v.GetString("postgres.host"),
			Port:         v.GetInt("postgres.port"),
			Database:     v.GetString("postgres.database"),
			Username:     v.GetString("postgres.username"),
			Password:     v.GetString("postgres.password"),
			SSLMode:      v.GetString("postgres.sslMode"),
			MaxOpenConns: v.GetInt("postgres.maxOpenConns"),
		},

		Log: LogConfig{
			Level:  v.GetString("logger.level"),
			Format: v.GetString("logger.format"),
			Output: v.GetString("logger.output"),
		},
	}
}
