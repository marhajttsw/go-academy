package db

import "project/entity"

func (db *MemoryDB) AddCharacter(movieTitle string, character entity.Character) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	id, exists := db.titleToID[movieTitle]
	if !exists {
		return false
	}
	character.MovieID = id
	db.characters[id] = append(db.characters[id], character)
	return true
}

func (db *MemoryDB) GetCharacters(movieTitle string) []entity.Character {
	db.mu.RLock()
	defer db.mu.RUnlock()
	id, exists := db.titleToID[movieTitle]
	if !exists {
		return nil
	}
	return db.characters[id]
}

func (db *MemoryDB) UpdateCharacter(movieTitle, oldName, newName string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	id, exists := db.titleToID[movieTitle]
	if !exists {
		return false
	}
	if chars, ok := db.characters[id]; ok {
		for i, char := range chars {
			if char.Name == oldName {
				db.characters[id][i].Name = newName
				return true
			}
		}
	}
	return false
}

func (db *MemoryDB) DeleteCharacter(movieTitle, name string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	id, exists := db.titleToID[movieTitle]
	if !exists {
		return false
	}
	if chars, ok := db.characters[id]; ok {
		for i, char := range chars {
			if char.Name == name {
				db.characters[id] = append(chars[:i], chars[i+1:]...)
				return true
			}
		}
	}
	return false
}
