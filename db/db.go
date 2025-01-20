package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var conn *pgx.Conn
var err error

func Init() {
	// Load environment variables
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load env: %v\n", err)
		os.Exit(1)
	}

	// Establish a connection to the PostgreSQL database
	conn, err = pgx.Connect(context.Background(), os.Getenv("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
}

func GetConn() *pgx.Conn {
	if conn == nil {
		fmt.Fprintf(os.Stderr, "Error: connection to database is not established")
		os.Exit(1)
	}
	return conn
}
