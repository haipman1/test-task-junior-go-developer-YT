package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`

	// Server
	ServerPort   string        `mapstructure:"SERVER_PORT"`
	ReadTimeout  time.Duration `mapstructure:"SERVER_READ_TIMEOUT"`
	WriteTimeout time.Duration `mapstructure:"SERVER_WRITE_TIMEOUT"`

	// App
	Environment string `mapstructure:"APP_ENV"`
	LogLevel    string `mapstructure:"LOG_LEVEL"`
}

func LoadConfig() (config Config, err error) {
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "1235")
	viper.SetDefault("DB_NAME", "tasks")
	viper.SetDefault("SERVER_PORT", ":8080")

	// Читаем из .env файла если есть
	viper.SetConfigName("config.env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("SERVER_PORT", ":8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("LOG_LEVEL", "debug")
	viper.SetDefault("SERVER_READ_TIMEOUT", "15s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "15s")

	viper.ReadInConfig()
	viper.AutomaticEnv()

	err = viper.Unmarshal(&config) //TODO error check here
	return
}
