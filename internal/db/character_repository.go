package db

import "project/internal/entity"

func (db *MemoryDB) AddCharacter(movieTitle string, character entity.Character) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	id, exists := db.titleToID[movieTitle]
	if !exists {
		return false
	}

	character.CharacterId = db.nextCharacterID
	db.nextCharacterID++
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

func (db *MemoryDB) GetCharacterByID(characterId uint64) (entity.Character, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	for _, chars := range db.characters {
		for _, c := range chars {
			if c.CharacterId == characterId {
				return c, true
			}
		}
	}
	return entity.Character{}, false
}

func (db *MemoryDB) UpdateCharacterByID(characterId uint64, updated entity.Character) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	for movieID, chars := range db.characters {
		for i, c := range chars {
			if c.CharacterId == characterId {
				updated.CharacterId = characterId
				if updated.MovieID == 0 {
					updated.MovieID = c.MovieID
				}
				db.characters[movieID][i] = updated
				return true
			}
		}
	}
	return false
}

func (db *MemoryDB) DeleteCharacterByID(characterId uint64) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	for movieID, chars := range db.characters {
		for i, c := range chars {
			if c.CharacterId == characterId {
				db.characters[movieID] = append(chars[:i], chars[i+1:]...)
				return true
			}
		}
	}
	return false
}
