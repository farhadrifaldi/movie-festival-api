package apis

import (
	"context"
	"net/http"
	"strconv"

	"github.com/farhadrifaldi/movie-festival-api/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type MovieInput struct {
	Title       string `json:"title"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Genres      string `json:"genres"`
	Artists     string `json:"artists"`
	URL         string `json:"url"`
	Rating      int    `json:"rating"`
	ViewCount   int    `json:"view_count"`
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
	var err error

	conn := db.GetConn()

	pageStr := c.Query("page")
	limitStr := c.Query("limit")
	search := c.Query("search")

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

	var rows pgx.Rows
	if search != "" {
		rows, err = conn.Query(context.Background(), `SELECT 
			id, title, image, description, duration, genres, artists, url, view_count, rating 
			FROM movies
			WHERE title ILIKE $1
			OR description ILIKE $1
			OR artists ILIKE $1
			OR genres ILIKE $1
			LIMIT $2 
			OFFSET $3`, "%"+search+"%", limit, offset)
	} else {
		rows, err = conn.Query(context.Background(), `SELECT 
			id, title, image, description, duration, genres, artists, url, view_count, rating 
			FROM movies
			LIMIT $1 
			OFFSET $2`, limit, offset)
	}
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

// Create movie function
func CreateMovie(c *gin.Context) {
	conn := db.GetConn()
	var movie MovieInput
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `INSERT INTO movies (
		title,
		image,
		description,
		duration,
		genres,
		artists,
		url,
		rating,
		view_count
	) VALUES ($1, $2, $3,$4, $5, $6, $7, $8, $9)`

	// Insert movie into the database
	_, err := conn.Exec(context.Background(), query, movie.Title, movie.Image, movie.Description, movie.Duration, movie.Genres, movie.Artists, movie.URL, movie.Rating, movie.ViewCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, movie)
}

// Update movie by ID
func UpdateMovie(c *gin.Context) {
	conn := db.GetConn()
	id := c.Param("id")
	var movie MovieInput
	if err := c.ShouldBindJSON(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE movies SET 
		title = $1,
		image = $2,
		description = $3,
		duration = $4,
		genres = $5,
		artists = $6,
		url = $7,
		rating = $8,
		view_count = $9
	WHERE id = $10`

	// Update movie in the database
	_, err := conn.Exec(context.Background(), query, movie.Title, movie.Image, movie.Description, movie.Duration, movie.Genres, movie.Artists, movie.URL, movie.Rating, movie.ViewCount, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, movie)
}

// Get movie by ID
func GetMovieByID(c *gin.Context) {
	conn := db.GetConn()
	id := c.Param("id")

	query := `SELECT 
		id, title, image, description, duration, genres, artists, url, view_count, rating 
		FROM movies 
		WHERE id = $1`
	var movie MovieResponse
	err := conn.QueryRow(context.Background(), query, id).Scan(&movie.ID, &movie.Title, &movie.Image, &movie.Description, &movie.Duration, &movie.Genres, &movie.Artists, &movie.URL, &movie.ViewCount, &movie.Rating)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}
	c.JSON(http.StatusOK, movie)
}

// Delete movie by ID
func DeleteMovie(c *gin.Context) {
	conn := db.GetConn()

	id := c.Param("id")
	_, err := conn.Exec(context.Background(), "DELETE FROM movies WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie: " + err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
