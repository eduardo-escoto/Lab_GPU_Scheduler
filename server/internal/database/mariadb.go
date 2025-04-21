package database

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/go-sql-driver/mysql"
)

func Connect(dataSourceName string) (*sql.DB, error) {
    db, err := sql.Open("mysql", dataSourceName)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
        return nil, err
    }

    if err := db.Ping(); err != nil {
        log.Fatalf("Error connecting to database: %v", err)
        return nil, err
    }

    fmt.Println("Successfully connected to the database")
    return db, nil
}

func ConnectMariaDB(dsn string) *sql.DB {
    db, err := sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to MariaDB: %v", err)
    }
    if err := db.Ping(); err != nil {
        log.Fatalf("MariaDB ping failed: %v", err)
    }
    log.Println("Connected to MariaDB successfully")
    return db
}