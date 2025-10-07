package db

import "project/entity"

func (db *MemoryDB) AddMovie(movie entity.Movie) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.movies[movie.ID] = movie
	db.titleToID[movie.Title] = movie.ID
}

func (db *MemoryDB) GetMovie(title string) (entity.Movie, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	id, exists := db.titleToID[title]
	if !exists {
		return entity.Movie{}, false
	}
	movie, ok := db.movies[id]
	return movie, ok
}

func (db *MemoryDB) UpdateMovie(title string, updated entity.Movie) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	id, exists := db.titleToID[title]
	if exists {
		db.movies[id] = updated
		return true
	}
	return false
}

func (db *MemoryDB) DeleteMovie(title string) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	id, exists := db.titleToID[title]
	if exists {
		delete(db.movies, id)
		delete(db.characters, id) // kinda cascade
		delete(db.titleToID, title)
		return true
	}
	return false
}

func (db *MemoryDB) ListMovies() []entity.Movie {
	db.mu.RLock()
	defer db.mu.RUnlock()
	movies := []entity.Movie{}
	for _, movie := range db.movies {
		movies = append(movies, movie)
	}
	return movies
}
