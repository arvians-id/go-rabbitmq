package setup

import (
	"database/sql"
	"fmt"
	"github.com/arvians-id/go-rabbitmq/gateway/cmd/config"
	_ "github.com/lib/pq"
	"time"
)

func TearDownTest(configuration config.Config) error {
	username := configuration.Get("DB_USERNAME_TEST")
	password := configuration.Get("DB_PASSWORD_TEST")
	host := configuration.Get("DB_HOST_TEST")
	port := configuration.Get("DB_PORT_TEST")
	database := configuration.Get("DB_DATABASE_TEST")
	sslMode := configuration.Get("DB_SSL_MODE_TEST")
	if configuration.Get("STATE") == "production" {
		username = configuration.Get("DB_USERNAME")
		password = configuration.Get("DB_PASSWORD")
		host = configuration.Get("DB_HOST")
		port = configuration.Get("DB_PORT")
		database = configuration.Get("DB_DATABASE")
		sslMode = configuration.Get("DB_SSL_MODE")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", host, port, username, password, database, sslMode)
	db, err := sql.Open(configuration.Get("DB_CONNECTION_TEST"), dsn)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(60)
	db.SetConnMaxLifetime(60 * time.Second)
	db.SetConnMaxIdleTime(60 * time.Second)

	_, err = db.Exec(`DELETE FROM todos;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM category_todos;`)
	if err != nil {
		return err
	}
	_, err = db.Exec(`DELETE FROM users;`)
	if err != nil {
		return err
	}

	return nil
}
