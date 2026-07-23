package bootstrap

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"go.uber.org/zap"
)

// Container holds all cross-cutting infrastructure dependencies.
type Container struct {
	Config *Config
	Logger *zap.Logger
	DB     *pgxpool.Pool
	Redis  *redis.Client
	Kafka  *kgo.Client
}

// BuildContainer initializes all infrastructure connections.
func BuildContainer(ctx context.Context) (*Container, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	logger, err := NewLogger(cfg.Log)
	if err != nil {
		return nil, err
	}

	db, err := NewDatabase(ctx, cfg.Database)
	if err != nil {
		logger.Error("Failed to connect to database", zap.Error(err))
		return nil, err
	}

	rdb, err := NewRedis(ctx, cfg.Redis)
	if err != nil {
		logger.Error("Failed to connect to redis", zap.Error(err))
		return nil, err
	}

	kafka, err := NewKafkaClient(cfg.Kafka)
	if err != nil {
		logger.Error("Failed to connect to kafka", zap.Error(err))
		return nil, err
	}

	return &Container{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Redis:  rdb,
		Kafka:  kafka,
	}, nil
}

func (c *Container) Close() {
	if c.Kafka != nil {
		c.Kafka.Close()
	}
	if c.Redis != nil {
		c.Redis.Close()
	}
	if c.DB != nil {
		c.DB.Close()
	}
	if c.Logger != nil {
		_ = c.Logger.Sync()
	}
}