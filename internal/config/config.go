package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config содержит все настройки приложения.
type Config struct {
	Database  DatabaseConfig  `mapstructure:"database"`
	Kafka     KafkaConfig     `mapstructure:"kafka"`
	GRPC      GRPCConfig      `mapstructure:"grpc"`
	Ingestion IngestionConfig `mapstructure:"ingestion"`
}

// DatabaseConfig содержит настройки подключения к базе данных.
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
}

// DSN возвращает строку подключения к базе данных.
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		d.Host, d.User, d.Password, d.DBName, d.Port, d.SSLMode)
}

// PostgresDSN возвращает строку подключения в формате postgres://.
func (d DatabaseConfig) PostgresDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.SSLMode)
}

// KafkaConfig содержит настройки Kafka.
type KafkaConfig struct {
	Brokers  []string            `mapstructure:"brokers"`
	Topic    string              `mapstructure:"topic"`
	Producer KafkaProducerConfig `mapstructure:"producer"`
}

// KafkaProducerConfig содержит настройки производителя Kafka.
type KafkaProducerConfig struct {
	BatchSize      int `mapstructure:"batch_size"`
	BatchTimeoutMs int `mapstructure:"batch_timeout_ms"`
	MaxAttempts    int `mapstructure:"max_attempts"`
	RequiredAcks   int `mapstructure:"required_acks"`
}

// GRPCConfig содержит настройки gRPC сервера.
type GRPCConfig struct {
	Port int `mapstructure:"port"`
}

// Address возвращает адрес для прослушивания gRPC сервера.
func (g GRPCConfig) Address() string {
	return fmt.Sprintf(":%d", g.Port)
}

// IngestionConfig содержит настройки сервиса ingestion.
type IngestionConfig struct {
	GHArchiveURL    string `mapstructure:"gharchive_url"`
	LookbackHours   int    `mapstructure:"lookback_hours"`
	Workers         int    `mapstructure:"workers"`
	ChannelSize     int    `mapstructure:"channel_size"`
	PollIntervalSec int    `mapstructure:"poll_interval_seconds"`
}

// LoadConfig загружает конфигурацию из файла.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}
