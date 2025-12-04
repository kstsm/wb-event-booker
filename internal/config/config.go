package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server    Server
	Postgres  Postgres
	Scheduler SchedulerConfig
	Telegram  TelegramConfig
}

type Server struct {
	Host string
	Port int
}

type Postgres struct {
	Username string
	Password string
	Host     string
	Port     string
	DBName   string
	SslMode  string
}

type SchedulerConfig struct {
	CheckInterval int
}

type TelegramConfig struct {
	BotToken string
}

func GetConfig() Config {
	viper.SetConfigFile(".env")

	err := viper.ReadInConfig()
	if err != nil {
		panic("Failed to read .env file: " + err.Error())
	}

	return Config{
		Server: Server{
			Host: viper.GetString("SRV_HOST"),
			Port: viper.GetInt("SRV_PORT"),
		},
		Postgres: Postgres{
			Username: viper.GetString("POSTGRES_USER"),
			Password: viper.GetString("POSTGRES_PASSWORD"),
			Host:     viper.GetString("POSTGRES_HOST"),
			Port:     viper.GetString("POSTGRES_PORT"),
			DBName:   viper.GetString("POSTGRES_DB"),
			SslMode:  viper.GetString("POSTGRES_SSLMODE"),
		},
		Scheduler: SchedulerConfig{
			CheckInterval: viper.GetInt("SCHEDULER_CHECK_INTERVAL"),
		},
		Telegram: TelegramConfig{
			BotToken: viper.GetString("TELEGRAM_BOT_TOKEN"),
		},
	}
}
