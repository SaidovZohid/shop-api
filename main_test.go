package main

import (
	"os"
	"testing"
	"fmt"
	"github.com/jmoiron/sqlx"
	"log"
)

var (
	dbManager *DBManager
)

func NewDbmanager(db *sqlx.DB) *DBManager {
	return &DBManager{db: db}
}

func TestMain(m *testing.M) {
	connStr := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatalf("failed to connect Database; %v", err)
	}
	dbManager = NewDbmanager(db)
 	os.Exit(m.Run())
}