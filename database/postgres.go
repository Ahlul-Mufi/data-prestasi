package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectDB() {
    dsn := "host=localhost port=5432 user=postgres password=your_password dbname=prestasi_db sslmode=disable"
    var err error
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("failed open db:", err)
    }
    if err = DB.Ping(); err != nil {
        log.Fatal("failed ping db:", err)
    }
    fmt.Println("Connected to prestasi_db ✅")
}
