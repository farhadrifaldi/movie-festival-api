package main

import (
	"context"
	"fmt"
	"os"

	"github.com/farhadrifaldi/movie-festival-api/apis"
	"github.com/farhadrifaldi/movie-festival-api/db"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func main() {
	var err error

	db.Init()

	conn = db.GetConn()

	r := gin.Default()

	// Define first API and make query from connected DB
	r.GET("/", sampleApifunc)

	// Define second API and make query from connected DB
	userMovieRoute := r.Group("/user/movies")
	{
		userMovieRoute.GET("/", apis.GetMovies)
		userMovieRoute.GET("/:id", apis.GetMovieByID)
	}

	// Define Router for admin
	adminRoute := r.Group("/admin")
	{
		// Movies API for admin
		adminRoute.GET("/movies", apis.GetMovies)
		adminRoute.POST("/movies", apis.CreateMovie)
		adminRoute.GET("/movies/:id", apis.GetMovieByID)
		adminRoute.PUT("/movies/:id", apis.UpdateMovie)
		adminRoute.DELETE("/movies/:id", apis.DeleteMovie)
	}

	r.Run() // listen and serve on 0.0.0.0:8080
}

func sampleApifunc(c *gin.Context) {
	var err error
	var title string
	err = conn.QueryRow(context.Background(), "SELECT title from movies WHERE id = $1", 1).Scan(&title)

	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	c.JSON(200, gin.H{
		"message": title,
	})
}
