package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/Kartik30R/sense/pkg/logger"
)

// Global Instances (Simplified Registry)
var (
	DB          *gorm.DB
	RedisClient *redis.Client
)

// Global Configuration Struct
type AppConfig struct {
	Env         string
	Port        string
	WorkerCount int
}

type DBConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

type KafkaConfig struct {
	Brokers      []string
	Topic        string
	GroupID      string
	BatchSize    int
	BatchTimeout time.Duration
	Async        bool
}

type RateLimitConfig struct {
	Enabled   bool
	GlobalRPS int
	DeviceRPS int
	BurstSize int
}

type FullConfig struct {
	App       AppConfig
	DB        DBConfig
	Redis     RedisConfig
	Kafka     KafkaConfig
	RateLimit RateLimitConfig
}

// LoadAllConfig loads values from environment variables
func LoadAllConfig() FullConfig {
	return FullConfig{
		App: AppConfig{
			Env:         getEnv("APP_ENV", "development"),
			Port:        getEnv("PORT", "8080"),
			WorkerCount: getEnvInt("WORKER_COUNT", 16),
		},
		DB: DBConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnvInt("DB_PORT", 5432),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "sense"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME_SEC", 3600)) * time.Second,
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnvInt("REDIS_PORT", 6379),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
			PoolSize: getEnvInt("REDIS_POOL_SIZE", 100),
		},
		Kafka: KafkaConfig{
			Brokers:      []string{getEnv("KAFKA_BROKERS", "localhost:9092")},
			Topic:        getEnv("KAFKA_TOPIC", "sensor_packets"),
			GroupID:      getEnv("KAFKA_GROUP_ID", "sensor_processor_group"),
			BatchSize:    getEnvInt("KAFKA_BATCH_SIZE", 100),
			BatchTimeout: time.Duration(getEnvInt("KAFKA_BATCH_TIMEOUT_MS", 10)) * time.Millisecond,
			Async:        getEnvBool("KAFKA_ASYNC", true),
		},
		RateLimit: RateLimitConfig{
			Enabled:   getEnvBool("RATE_LIMIT_ENABLED", true),
			GlobalRPS: getEnvInt("RATE_LIMIT_GLOBAL_RPS", 100000),
			DeviceRPS: getEnvInt("RATE_LIMIT_DEVICE_RPS", 50),
			BurstSize: getEnvInt("RATE_LIMIT_BURST_SIZE", 100),
		},
	}
}

// InitDB initializes PostgreSQL connection
func InitDB(cfg DBConfig) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:          gormlogger.Default.LogMode(gormlogger.Info),
		CreateBatchSize: 1000,
	})

	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)
}

// InitRedis initializes Redis connection
func InitRedis(cfg RedisConfig) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.PoolSize / 10,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := RedisClient.Ping(ctx).Result(); err != nil {
		logger.Fatal("Redis connection failed", zap.Error(err))
	}
}

// Helpers
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(value); err == nil {
			return b
		}
	}
	return fallback
}
