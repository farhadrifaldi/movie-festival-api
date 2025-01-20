package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

var conn *pgx.Conn

func main() {
	var err error
	// Establish a connection to the PostgreSQL database
	conn, err = pgx.Connect(context.Background(), goDotEnvVariable("POSTGRES_CONNECTION_STRING"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	r := gin.Default()

	// Define first API and make query from connected DB
	r.GET("/", func(c *gin.Context) {
		var title string
		err = conn.QueryRow(context.Background(), "SELECT title from movies WHERE id = $1", 1).Scan(&title)

		if err != nil {
			fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
			os.Exit(1)
		}

		c.JSON(200, gin.H{
			"message": title,
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

/*
*
Function to get environment variable from .env
*/
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot load env")
		os.Exit(1)
	}

	return os.Getenv(key)
}
