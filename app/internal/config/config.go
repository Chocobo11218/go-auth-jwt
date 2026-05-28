package config

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
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	Expire string `mapstructure:"expire"`
}

type SecretConfig struct {
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	JWTSecret  string `envconfig:"JWT_SECRET" required:"true"`
}
