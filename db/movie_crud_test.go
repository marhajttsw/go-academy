package db

import (
	"project/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMovieCRUD(t *testing.T) {
	db := New()
	movie := entity.Movie{ID: 1, Title: "Inception", Year: 2010}

	// AddMovie
	db.AddMovie(movie)
	got, ok := db.GetMovie("Inception")
	assert.True(t, ok, "AddMovie/GetMovie failed")
	assert.Equal(t, "Inception", got.Title)

	// UpdateMovie
	updated := entity.Movie{ID: 1, Title: "Inception", Year: 2011}
	assert.True(t, db.UpdateMovie("Inception", updated), "UpdateMovie failed")
	got, _ = db.GetMovie("Inception")
	assert.Equal(t, 2011, got.Year, "UpdateMovie did not update year")

	// ListMovies
	movies := db.ListMovies()
	assert.Len(t, movies, 1, "ListMovies failed")

	// DeleteMovie
	assert.True(t, db.DeleteMovie("Inception"), "DeleteMovie failed")
	_, ok = db.GetMovie("Inception")
	assert.False(t, ok, "DeleteMovie did not remove movie")
}
