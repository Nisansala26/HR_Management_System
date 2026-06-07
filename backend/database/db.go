package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "erp_user")
	password := getEnv("DB_PASSWORD", "erp_secure_password")
	dbname := getEnv("DB_NAME", "erp_db")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	for i := 0; i < 10; i++ {
		log.Printf("Connecting to database (attempt %d/10)...", i+1)
		DB, err = sql.Open("postgres", connStr)
		if err == nil {
			err = DB.Ping()
			if err == nil {
				log.Println("Successfully connected to the database!")
				return
			}
		}
		log.Printf("Failed to connect: %v. Retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}

	log.Fatalf("Could not connect to database: %v", err)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
