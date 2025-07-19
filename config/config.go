package config

import "github.com/spf13/viper"

var Configs *Config

type Config struct {
	Env           string `mapstructure:"APP_ENV"`
	AppPort       int    `mapstructure:"APP_PORT"`
	AppHost       string `mapstructure:"APP_HOST"`
	AppName       string `mapstructure:"APP_NAME"`
	PgUri         string `mapstructure:"POSTGRES_URI"`
	RedisAddress  string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	SmtpHost      string `mapstructure:"SMTP_HOST"`
	SmtpPort      int    `mapstructure:"SMTP_PORT"`
	SmtpUser      string `mapstructure:"SMTP_USER"`
	SmtpPassword  string `mapstructure:"SMTP_PASSWORD"`
	SmtpFrom      string `mapstructure:"SMTP_FROM"`
	JWTSecret     string `mapstructure:"JWT_SECRET"`
	HolidayApiKey string `mapstructure:"HOLIDAY_API_KEY"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&Configs)
	return Configs, err
}
