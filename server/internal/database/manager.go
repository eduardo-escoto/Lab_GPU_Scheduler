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

// QueryAndMap executes a query and maps the rows to a slice of a specific type using a mapper function.
func QueryAndMap[T any](db *sql.DB, query string, args []interface{}, mapper func(*sql.Rows) (T, error)) ([]T, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []T
	for rows.Next() {
		item, err := mapper(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
