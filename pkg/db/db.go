// Package db defines functionality for connecting to a relational database using gorm.
package db

import (
	"context"
	"maker-checker/config"
	"math/rand"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB represents the database instance for the whole app.
type DB struct {
	db       *gorm.DB
	readDSN  string
	writeDSN string
}

// Init returns a new instance of database object with connection.
func Init(ctx context.Context, cfg *config.AppConfig) (*DB, error) {
	rand.Seed(time.Now().UnixNano())

	db, err := gorm.Open(postgres.Open(cfg.Database.WriteDSN), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, err
	}

	return &DB{
		db:       db.WithContext(ctx),
		readDSN:  cfg.Database.ReadDSN,
		writeDSN: cfg.Database.WriteDSN,
	}, nil
}

// Stop stops the database connection.
func (db *DB) Stop() error {
	if db.db == nil {
		return nil
	}

	pqDB, err := db.db.DB()
	if err != nil {
		return err
	}

	return pqDB.Close()
}

// DB returns the gorm db instance.
func (db *DB) DB() *gorm.DB {
	return db.db
}
