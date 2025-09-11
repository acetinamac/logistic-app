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

	// Ensure ALL PostgreSQL enum types exist
	if err := database.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN CREATE TYPE user_role_enum AS ENUM ('client','admin'); END IF; END $$;").Error; err != nil {
		return nil, err
	}

	if err := database.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'package_size_enum') THEN CREATE TYPE package_size_enum AS ENUM ('S','M','L','XL'); END IF; END $$;").Error; err != nil {
		return nil, err
	}

	if err := database.Exec("DO $$ BEGIN IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'order_status_enum') THEN CREATE TYPE order_status_enum AS ENUM ('created','collected','in_station','in_route','delivered','cancelled'); END IF; END $$;").Error; err != nil {
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
