// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package db

import (
	"context"
	"database/sql"
	"log"
	"log/slog"

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
	POSTGRES_DB_USER_PASSWORD := os.Getenv("DB_PASS")
	POSTGRES_HOST := os.Getenv("DB_HOST")
	POSTGRES_DB_NAME := os.Getenv("DB_NAME")
	POSTGRES_DB_USER := os.Getenv("DB_USER")
	dbURI := fmt.Sprintf("host=%s user=%s password=%s port=%s database=%s",
		POSTGRES_HOST, POSTGRES_DB_USER, POSTGRES_DB_USER_PASSWORD, "5432", POSTGRES_DB_NAME)
	log.Println(dbURI)

	db, err := sql.Open("pgx", dbURI)

	if err != nil {
		return nil, fmt.Errorf("Error open sql url: %v", err)
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(2)

	err = db.PingContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Error pinging database: %v", err)
	}
	slog.Log(context.Background(), slog.LevelInfo, "DB pinged successfully")
	return db, nil
}
