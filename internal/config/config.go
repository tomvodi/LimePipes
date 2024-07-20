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
