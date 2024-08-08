package config

type Config struct {
	ServerUrl string `mapstructure:"API_SERVER_URL"`

	TlsCertPath    string `mapstructure:"TLS_CERT_PATH"`
	TlsCertKeyPath string `mapstructure:"TLS_CERT_KEY_PATH"`

	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     string `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbSslMode  string `mapstructure:"DB_SSL_MODE"`
	DbTimeZone string `mapstructure:"DB_TIMEZONE"`

	HealthCacheDurationSeconds uint32 `mapstructure:"HEALTH_CACHE_DURATION_SECONDS"`
	HealthGlobalTimeoutSeconds uint32 `mapstructure:"HEALTH_GLOBAL_TIMEOUT_SECONDS"`
	HealthRefreshPeriodSeconds uint32 `mapstructure:"HEALTH_REFRESH_PERIOD_SECONDS"`
	HealthInitialDelaySeconds  uint32 `mapstructure:"HEALTH_INITIAL_DELAY_SECONDS"`

	PluginsDirectoryPath string `mapstructure:"PLUGINS_DIRECTORY_PATH"`
}

func (c *Config) DbConfig() DbConfig {
	return DbConfig{
		Host:     c.DbHost,
		Port:     c.DbPort,
		DbName:   c.DbName,
		User:     c.DbUser,
		Password: c.DbPassword,
		SslMode:  c.DbSslMode,
		TimeZone: c.DbTimeZone,
	}
}

func (c *Config) HealthConfig() HealthConfig {
	return HealthConfig{
		CacheDurationSeconds: c.HealthCacheDurationSeconds,
		GlobalTimeoutSeconds: c.HealthGlobalTimeoutSeconds,
		RefreshPeriodSeconds: c.HealthRefreshPeriodSeconds,
		InitialDelaySeconds:  c.HealthInitialDelaySeconds,
	}
}
