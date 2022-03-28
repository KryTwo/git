package postgres

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func InitDB() {
	var err error
	connStr := "user=root password=123456 dbname=postgres sslmode=disable"
	Db, err = sql.Open("postgres", connStr)
	Db.SetMaxIdleConns(100)
	if err != nil {
		log.Fatal(err)
	}
}
