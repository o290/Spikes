package config

type Config struct {
	Mysql  Mysql
	System System
	Logger Logger
	Redis  Redis
	JWT    Jwt
}
