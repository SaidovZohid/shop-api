package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	host = "localhost"
	port = 5432
	user = "postgres"
	password = "1234"
	dbname = "shop"
)

type DBManager struct {
	db *sqlx.DB
}

func NewDBmanager(db *sqlx.DB) *DBManager {
	return &DBManager{db: db}
}

var (
	Dbmanager *DBManager
)

func main(){
	connStr := fmt.Sprintf("host = %s port = %d user = %s password = %s dbname = %s sslmode=disable", host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", connStr)

	if err != nil {
		log.Fatalf("failed to connect Database; %v", err)
	}

	fmt.Println("Succesfully Connected!")
	Dbmanager = NewDBmanager(db)

	server := NewServer(Dbmanager)

	err = server.Run(":8080")
	if err != nil {
		log.Fatalf("failed to start server; %v", err)
	}
}