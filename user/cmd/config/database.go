package config

import (
	"context"
	"database/sql"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

func NewInitializedDatabase(configuration Config) (*mongo.Database, error) {
	db, err := NewMongoNoSQL(configuration)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewMongoNoSQL(configuration Config) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := configuration.Get("DB_USERNAME")
	password := configuration.Get("DB_PASSWORD")
	host := configuration.Get("DB_HOST")
	port := configuration.Get("DB_PORT")
	database := configuration.Get("DB_DATABASE")

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s", username, password, host, port)
	clientOption := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOption)
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return client.Database(database), nil
}

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

	err = databasePooling(configuration, db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func databasePooling(configuration Config, db *sql.DB) error {
	// Limit connection with db pooling
	setMaxIdleConns, err := strconv.Atoi(configuration.Get("DATABASE_POOL_MIN"))
	if err != nil {
		return err
	}
	setMaxOpenConns, err := strconv.Atoi(configuration.Get("DATABASE_POOL_MAX"))
	if err != nil {
		return err
	}
	setConnMaxIdleTime, err := strconv.Atoi(configuration.Get("DATABASE_MAX_IDLE_TIME_SECOND"))
	if err != nil {
		return err
	}
	setConnMaxLifetime, err := strconv.Atoi(configuration.Get("DATABASE_MAX_LIFE_TIME_SECOND"))
	if err != nil {
		return err
	}

	db.SetMaxIdleConns(setMaxIdleConns)                                    // minimal connection
	db.SetMaxOpenConns(setMaxOpenConns)                                    // maximal connection
	db.SetConnMaxLifetime(time.Duration(setConnMaxIdleTime) * time.Second) // unused connections will be deleted
	db.SetConnMaxIdleTime(time.Duration(setConnMaxLifetime) * time.Second) // connection that can be used

	return nil
}
