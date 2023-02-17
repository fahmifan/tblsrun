package dbtool

import (
	"database/sql"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func OpenDB(dsn string) (*sql.DB, error) {
	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return gormDB.DB()
}

func CreateDB(db *sql.DB, dbName string) (err error) {
	if _, err = db.Exec(`CREATE DATABASE "` + dbName + `";`); err != nil {
		return err
	}
	return nil
}

func CreateSchemaIfNoExist(db *sql.DB, schema string) (err error) {
	// create & set the real schema
	if _, err = db.Exec(`CREATE SCHEMA IF NOT EXISTS "` + schema + `";`); err != nil {
		return err
	}
	return nil
}
