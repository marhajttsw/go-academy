package db

import (
	"project/entity"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCharacterCRUD(t *testing.T) {
	db := New()
	movie := entity.Movie{ID: 1, Title: "Inception", Year: 2010}
	db.AddMovie(movie)

	char := entity.Character{Movie: "Inception", Name: "Dom Cobb"}
	assert.True(t, db.AddCharacter("Inception", char), "AddCharacter failed")

	chars := db.GetCharacters("Inception")
	assert.Len(t, chars, 1, "GetCharacters failed")
	assert.Equal(t, "Dom Cobb", chars[0].Name)

	assert.True(t, db.UpdateCharacter("Inception", "Dom Cobb", "Arthur"), "UpdateCharacter failed")
	chars = db.GetCharacters("Inception")
	assert.Equal(t, "Arthur", chars[0].Name, "UpdateCharacter did not update name")

	assert.True(t, db.DeleteCharacter("Inception", "Arthur"), "DeleteCharacter failed")
	chars = db.GetCharacters("Inception")
	assert.Len(t, chars, 0, "DeleteCharacter did not remove character")
}
