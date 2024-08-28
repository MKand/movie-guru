package db

import (
	"database/sql"
	"log"

	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type MovieDB struct {
	DB *sql.DB
}

func GetDB() (*MovieDB, error) {
	db, err := connectToDB()
	if err != nil {
		return nil, err
	}

	return &MovieDB{DB: db}, nil
}

func connectToDB() (*sql.DB, error) {
	POSTGRES_DB_USER_PASSWORD := os.Getenv("POSTGRES_DB_USER_PASSWORD")
	POSTGRES_HOST := os.Getenv("POSTGRES_HOST")
	POSTGRES_DB_NAME := os.Getenv("POSTGRES_DB_NAME")
	POSTGRES_DB_USER := os.Getenv("POSTGRES_DB_USER")
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		POSTGRES_HOST, POSTGRES_DB_USER, POSTGRES_DB_USER_PASSWORD, "5432", POSTGRES_DB_NAME)
	log.Println(dbURI)
	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
