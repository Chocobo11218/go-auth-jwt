package config

import "time"

type AppConfig struct {
	Log    Log          `mapstructure:"log" validate:"required"`
	App    App          `mapstructure:"app" validate:"required"`
	Server Server       `mapstructure:"server" validate:"required"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Secret SecretConfig `mapstructure:"-"`
}

type Log struct {
	Level string `mapstructure:"level" validate:"required"`
	Env   string `mapstructure:"env" validate:"required"`
}

type App struct {
	Name      string `mapstructure:"name" validate:"required"`
	ProjectID string `mapstructure:"project-id" validate:"required"`
}

type Server struct {
	Address  string `mapstructure:"address" validate:"required"`
	TimeZone string `mapstructure:"time-zone" validate:"required"`
}

type DBConfig struct {
	Name                      string        `mapstructure:"name"`
	SSLMode                   string        `mapstructure:"tls" validate:"required"`
	MaxOpenConnection         int           `mapstructure:"max_open_conns" validate:"required"`
	MaxIdleConnection         int           `mapstructure:"max_idle_conns" validate:"required"`
	ConnectionMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" validate:"required"`
}

type JWTConfig struct {
	Expire time.Duration `mapstructure:"expire"`
}

type SecretConfig struct {
	DBHost     string `envconfig:"DB_HOST" required:"true"`
	DBPort     string `envconfig:"DB_PORT" required:"true"`
	DBUser     string `envconfig:"DB_USER" required:"true"`
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	JWTSecret  string `envconfig:"JWT_SECRET" required:"true"`
}
