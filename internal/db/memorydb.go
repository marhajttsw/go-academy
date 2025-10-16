package db

import (
	"project/internal/entity"
	"sync"
)

type MemoryDB struct {
	movies     map[uint64]entity.Movie
	characters map[uint64][]entity.Character
	titleToID  map[string]uint64
	nextCharacterID uint64
	nextMovieID     uint64
	mu         sync.RWMutex
}

func New() *MemoryDB {
	return &MemoryDB{
		movies:     make(map[uint64]entity.Movie),
		characters: make(map[uint64][]entity.Character),
		titleToID:  make(map[string]uint64),
		nextCharacterID: 1,
		nextMovieID:     1,
	}
}
