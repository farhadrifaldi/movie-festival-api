package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"

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

	// Movies API
	r.GET("/movies", GetMovies)
	r.GET("/movies/:id", GetMovieByID)

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

type MovieResponse struct {
	ID          string
	Title       string
	Image       string
	Description string
	Duration    int // in seconds
	Genres      string
	Artists     string
	URL         string
	Rating      int
	ViewCount   int
}

// Get movie function with page & limit as parameter
func GetMovies(c *gin.Context) {
	pageStr := c.Query("page")
	limitStr := c.Query("limit")

	// Convert page and limit to integers
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1 // default to page 1
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10 // default limit
	}

	offset := (page - 1) * limit

	rows, err := conn.Query(context.Background(), "SELECT id, title, image, description, duration, genres, artists, url, view_count, rating FROM movies LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to execute query: " + err.Error()})
		return
	}
	defer rows.Close()

	var movies []MovieResponse
	for rows.Next() {
		var movie MovieResponse
		if err = rows.Scan(&movie.ID,
			&movie.Title,
			&movie.Image,
			&movie.Description,
			&movie.Duration,
			&movie.Genres,
			&movie.Artists,
			&movie.URL,
			&movie.ViewCount,
			&movie.Rating); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan row: " + err.Error()})
			return
		}
		movies = append(movies, movie)
	}
	c.JSON(http.StatusOK, movies)
}

// Get movie by ID
func GetMovieByID(c *gin.Context) {
	id := c.Param("id")
	var movie MovieResponse
	err := conn.QueryRow(context.Background(), "SELECT id, title, image, description, duration, genres, artists, url, view_count, rating FROM movies WHERE id = $1", id).Scan(&movie.ID, &movie.Title, &movie.Image, &movie.Description, &movie.Duration, &movie.Genres, &movie.Artists, &movie.URL, &movie.ViewCount, &movie.Rating)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	c.JSON(http.StatusOK, movie)
}
