package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx" // SQL library that simplifies database work
	_ "github.com/lib/pq"     // PostgreSQL driver, imported for side effects only
)

// DB is a package-level variable holding the database connection pool
var DB *sqlx.DB

// Connect initializes the connection to the PostgreSQL database
func Connect() {

	dsn := "user=himanshu password=himu dbname=himu sslmode=disable"

	// Connect to PostgreSQL using sqlx and the DSN (data source name)
	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	fmt.Println("✅ Connected to PostgreSQL database!")
}
