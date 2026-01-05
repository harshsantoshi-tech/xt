package configs

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var AppConfig struct {
	ServerPort string
	DB         *sql.DB
}

func init() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, using system environment variables")
	}

	// Set application-wide configurations
	AppConfig = struct {
		ServerPort string
		DB         *sql.DB
	}{
		ServerPort: getPort(),
		DB:         ConnectDB(),
	}
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
