package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct{ *gorm.DB }

func Connect() (*Database, error) {
	host := getenv("POSTGRES_HOST", "localhost")
	port := getenv("POSTGRES_PORT", "5432")
	user := getenv("POSTGRES_USER", "postgres")
	pass := getenv("POSTGRES_PASSWORD", "postgres")
	dbname := getenv("POSTGRES_DB", "logistics")
	ssl := getenv("POSTGRES_SSLMODE", "disable")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC", host, user, pass, dbname, port, ssl)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	log.Println("connected to postgres")
	return &Database{database}, nil
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
