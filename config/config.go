package config

type Config struct {
	App      AppConfig
	LDAP     LDAPConfig
	Postgres PostgresConfig
}

type AppConfig struct {
	Name          string
	Env           string
	Port          int
	Debug         bool
	Timezone      string
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpiry     string `mapstructure:"jwt_expiry"`
	RefreshExpiry string `mapstructure:"refresh_expiry"`
}

type LDAPConfig struct {
	Host     string
	Port     int
	BaseDN   string
	BindDN   string
	BindPass string
	UseTLS   bool
}

type PostgresConfig struct {
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	SSLMode      string
	MaxOpenConns int
}
