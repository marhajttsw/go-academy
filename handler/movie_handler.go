package handler

import (
	"fmt"
	"project/db"
	"project/entity"
)

func RegisterMovieHandler(db *db.MemoryDB) {

	db.AddMovie(entity.Movie{Title: "Inception", Year: 2010})
	db.AddMovie(entity.Movie{Title: "The Matrix", Year: 1999})

	fmt.Println("Movies:")
	for _, movie := range db.ListMovies() {
		fmt.Printf(" %s %d\n", movie.Title, movie.Year)
	}
}
