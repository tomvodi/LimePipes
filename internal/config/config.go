package config

type Config struct {
	ServerUrl      string `mapstructure:"SERVER_URL"`
	SqliteDbPath   string `mapstructure:"SQLITE_DB_PATH"`
	TlsCertPath    string `mapstructure:"TLS_CERT_PATH"`
	TlsCertKeyPath string `mapstructure:"TLS_CERT_KEY_PATH"`
}
