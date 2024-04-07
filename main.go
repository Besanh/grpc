package main

import (
	"anhle-grpc/lib/cache"
	"anhle-grpc/lib/database"
	"anhle-grpc/lib/dotenv"
	"anhle-grpc/lib/elasticsearch"
	"anhle-grpc/lib/redis"
	"anhle-grpc/repository"
	"io"
	"os"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	// Database
	Host         string
	Port         int
	Database     string
	Username     string
	Password     string
	Timeout      int
	ReadTimeout  int
	WriteTimeout int
	DialTimeout  int
	MaxIdleConns int
	MaxOpenConns int

	// Log
	LogLevel string
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	config := database.PostgreSql{
		Host:         dotenv.GetStringENV("DB_HOST", "localhost"),
		Port:         dotenv.GetIntENV("DB_PORT", 5432),
		Database:     dotenv.GetStringENV("DB_DATABASE", "postgres"),
		Username:     dotenv.GetStringENV("DB_USERNAME", "postgres"),
		Password:     dotenv.GetStringENV("DB_PASSWORD", "postgres"),
		Timeout:      dotenv.GetIntENV("DB_TIMEOUT", 30),
		ReadTimeout:  dotenv.GetIntENV("DB_READ_TIMEOUT", 30),
		WriteTimeout: dotenv.GetIntENV("DB_WRITE_TIMEOUT", 30),
		DialTimeout:  dotenv.GetIntENV("DB_DIAL_TIMEOUT", 20),
		MaxIdleConns: dotenv.GetIntENV("DB_MAX_IDLE_CONNS", 20),
		MaxOpenConns: dotenv.GetIntENV("DB_MAX_OPEN_CONNS", 10),
	}
	repository.DBCON = database.NewPostgreSqlCon(config)

	esConfig := elasticsearch.Config{
		Username:      dotenv.GetStringENV("ES_USERNAME", ""),
		Password:      dotenv.GetStringENV("ES_PASSWORD", ""),
		Host:          []string{dotenv.GetStringENV("ES_HOST", "")},
		RetryStatuses: []int{502, 503, 504},
		MaxRetries:    dotenv.GetIntENV("ES_MAX_RETRIES", 3),
	}

	repository.ESCON = elasticsearch.NewElasticsearch(esConfig)

	redisConfig := redis.Config{
		Addr:         dotenv.GetStringENV("REDIS_ADDR", "localhost:6379"),
		Password:     dotenv.GetStringENV("REDIS_PASSWORD", ""),
		DB:           dotenv.GetIntENV("REDIS_DB", 0),
		PoolSize:     dotenv.GetIntENV("REDIS_POOL_SIZE", 10),
		PoolTimeout:  dotenv.GetIntENV("REDIS_POOL_TIMEOUT", 10),
		IdleTimeout:  dotenv.GetIntENV("REDIS_IDLE_TIMEOUT", 10),
		ReadTimeout:  dotenv.GetIntENV("REDIS_READ_TIMEOUT", 10),
		WriteTimeout: dotenv.GetIntENV("REDIS_WRITE_TIMEOUT", 10),
	}
	var err error
	redis.Redis, err = redis.NewRedis(redisConfig)
	if err != nil {
		panic(err)
	}
}

func main() {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer file.Close()
	setAppLogger(Config{
		LogLevel: "info",
	}, file)

	cache.RCache = cache.NewRedisCache(redis.Redis.GetClient())
}

func setAppLogger(cfg Config, file *os.File) {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: false,
	})
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})
	switch cfg.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(io.MultiWriter(os.Stdout, file))
}
