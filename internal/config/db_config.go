package config

type DbConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	SslMode  string
	TimeZone string
}

type HealthConfig struct {
	CacheDurationSeconds uint32
	GlobalTimeoutSeconds uint32
	RefreshPeriodSeconds uint32
	InitialDelaySeconds  uint32
}
