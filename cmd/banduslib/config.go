package main

type Config struct {
	ServerUrl    string `mapstructure:"SERVER_URL"`
	SqliteDbPath string `mapstructure:"SQLITE_DB_PATH"`
}
