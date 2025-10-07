package handler

import (
	"fmt"
	"project/db"
	"project/entity"
)

func RegisterCharacterHandler(db *db.MemoryDB) {
	db.AddCharacter("Inception", entity.Character{Movie: "Inception", Name: "Dom Cobb"})
	db.AddCharacter("Inception", entity.Character{Movie: "Inception", Name: "Arthur"})
	db.AddCharacter("The Matrix", entity.Character{Movie: "The Matrix", Name: "Neo"})
	db.AddCharacter("The Matrix", entity.Character{Movie: "The Matrix", Name: "Morpheus"})

	updated := db.UpdateCharacter("Inception", "Arthur", "Eames")
	fmt.Printf("\nUpdate 'Arthur' to 'Eames' in 'Inception': %v\n", updated)

	deleted := db.DeleteCharacter("The Matrix", "Neo")
	fmt.Printf("Delete 'Neo' from 'The Matrix': %v\n", deleted)

	fmt.Println("\nCharacters in 'Inception':")
	for _, char := range db.GetCharacters("Inception") {
		fmt.Printf("%s\n", char.Name)
	}

	fmt.Println("\nCharacters in 'The Matrix':")
	for _, char := range db.GetCharacters("The Matrix") {
		fmt.Printf("%s\n", char.Name)
	}
}
