package bootstrap

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Log      LogConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port int
}

type LogConfig struct {
	Level  string
	Format string
}

type DatabaseConfig struct {
	URL      string
	MaxConns int32 `mapstructure:"max_conns"`
}

type RedisConfig struct {
	URL string
}

type KafkaConfig struct {
	Brokers []string
}

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")

	// Environment variable overrides (e.g., GHOSTKEY_DATABASE_URL)
	v.SetEnvPrefix("GHOSTKEY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}