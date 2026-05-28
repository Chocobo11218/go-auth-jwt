package config

type Config struct {
	App    AppConfig    `mapstructure:"app"`
	DB     DBConfig     `mapstructure:"db"`
	JWT    JWTConfig    `mapstructure:"jwt"`
	Secret SecretConfig `mapstructure:"-"`
}

type AppConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type DBConfig struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Name string `mapstructure:"name"`
}

type JWTConfig struct {
	Expire string `mapstructure:"expire"`
}

type SecretConfig struct {
	DBPassword string `envconfig:"DB_PASSWORD" required:"true"`
	JWTSecret  string `envconfig:"JWT_SECRET" required:"true"`
}
