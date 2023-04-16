package config

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func NewPostgresSQL(configuration Config) (*sql.DB, error) {
	username := configuration.Get("DB_USERNAME")
	password := configuration.Get("DB_PASSWORD")
	host := configuration.Get("DB_HOST")
	port := configuration.Get("DB_PORT")
	database := configuration.Get("DB_DATABASE")
	sslMode := configuration.Get("DB_SSL_MODE")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, database, sslMode)
	db, err := sql.Open(configuration.Get("DB_CONNECTION"), dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db, err = databasePooling(configuration, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func databasePooling(configuration Config, db *sql.DB) (*sql.DB, error) {
	setMaxIdleConns, err := strconv.Atoi(configuration.Get("DB_POOL_MIN"))
	if err != nil {
		return nil, err
	}
	setMaxOpenConns, err := strconv.Atoi(configuration.Get("DB_POOL_MAX"))
	if err != nil {
		return nil, err
	}
	setConnMaxIdleTime, err := strconv.Atoi(configuration.Get("DB_MAX_IDLE_TIME"))
	if err != nil {
		return nil, err
	}
	setConnMaxLifetime, err := strconv.Atoi(configuration.Get("DB_MAX_LIFE_TIME"))
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(setMaxIdleConns)                                    // minimal connection
	db.SetMaxOpenConns(setMaxOpenConns)                                    // maximal connection
	db.SetConnMaxLifetime(time.Duration(setConnMaxIdleTime) * time.Second) // unused connections will be deleted
	db.SetConnMaxIdleTime(time.Duration(setConnMaxLifetime) * time.Second)

	return db, nil
}
