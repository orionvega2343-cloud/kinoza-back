package config

import (
	"log"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Db     DbConfig     `yaml:"db"`
	Server ServerConfig `yaml:"server"`
	Redis  Redis        `yaml:"redis"`
	Kafka  Kafka        `yaml:"kafka"`
	Jwt    Jwt          `yaml:"jwt"`
}

type DbConfig struct {
	Name     string `yaml:"name" env-required:"true"`
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `env:"DB_PASS" env-required:"true"`
	SslMode  string `yaml:"ssl_mode" env-required:"true"`
}

type ServerConfig struct {
	Host string `yaml:"host" env-required:"true"`
	Port int    `yaml:"port" env-required:"true"`
}

type Redis struct {
	Addr     string `yaml:"addr" env-required:"true"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

type Kafka struct {
	Brokers       []string `yaml:"brokers" env-required:"true"`
	ConsumerGroup string   `yaml:"consumer_group" env-required:"true"`
}

type Jwt struct {
	Secret string `env:"JWT_KEY" env-required:"true"`
}

// MustLoad - получает данные из .env, вызывает собранный конфиг,
// загружает его через чтение .yaml файла
func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		slog.Error("failed to load .env file", "error", err)
	}

	cfg := Config{}
	err = cleanenv.ReadConfig("configs/config_dev.yml", &cfg)
	if err != nil {
		log.Fatalf("failed to load config file %v", err)
	}
	return &cfg
}
