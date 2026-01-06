package configs

import (
	"context" // Required for Redis
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-redis/redis/v8" // Import Redis
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var AppConfig struct {
	ServerPort string
	DB         *sql.DB
	Redis      *redis.Client // Added Redis Client
}

var ctx = context.Background()

func InitializeConfigs() {
	// 1. Load Env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	AppConfig.ServerPort = getPort()
	AppConfig.DB = ConnectDB()

	AppConfig.Redis = ConnectRedis()
}

func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	return ":" + port
}

func ConnectDB() *sql.DB {
	dbUser := os.Getenv("DB_USERNAME")
	dbPass := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	if dbUser == "" || dbPass == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("❌ Database environment variables not fully set")
	}

	// DSN format: user:password@tcp(host:port)/dbname
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&tls=true",
		dbUser, dbPass, dbHost, dbPort, dbName,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to open DB: %v", err)
	}

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(time.Minute * 3)

	if err = db.Ping(); err != nil {
		log.Fatalf("❌ DB not reachable: %v", err)
	}

	fmt.Println("✅ Database connected")
	return db
}

// ConnectRedis sets up the Redis connection
func ConnectRedis() *redis.Client {
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPass := os.Getenv("REDIS_PASSWORD")

	if redisHost == "" || redisPort == "" {
		log.Println("⚠️ Redis environment variables not set, defaulting to localhost:6379")
		redisHost = "localhost"
		redisPort = "6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPass,
		DB:       0,
	})

	// Ping Redis to check if it's reachable
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("❌ Redis not reachable: %v", err)
	}

	fmt.Println("✅ Redis connected")
	return rdb
}
