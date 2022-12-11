package test

import (
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"simple-open-library/helper"
	"time"

	"github.com/joho/godotenv"
)

func SetupTestDB() *sql.DB {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	err := godotenv.Load(filepath.Join(basepath, "/../.env"))
	helper.FatalIfError(err, "Error loading .env file")

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_TEST_NAME")
	dbConn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"
	db, err := sql.Open("mysql", dbConn)
	helper.PanicIfError(err)

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(10 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)
	return db
}

func TruncateDatabase(db *sql.DB) {
	db.Exec("SET FOREIGN_KEY_CHECKS = 0")
	db.Exec("TRUNCATE pickup")
	db.Exec("TRUNCATE authored")
	db.Exec("TRUNCATE book")
	db.Exec("TRUNCATE author")
	db.Exec("SET FOREIGN_KEY_CHECKS = 1")
}
