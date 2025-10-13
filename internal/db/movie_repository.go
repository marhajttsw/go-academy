package db

import "project/internal/entity"

func (db *MemoryDB) AddMovie(movie entity.Movie) {
	db.mu.Lock()
	defer db.mu.Unlock()

	if movie.ID == 0 {
		movie.ID = db.nextMovieID
		db.nextMovieID++
	}
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
		// Preserve the original ID if caller didn't set it or set a different one.
		updated.ID = id
		// If title changed, update the title index mapping
		if updated.Title != title {
			delete(db.titleToID, title)
			db.titleToID[updated.Title] = id
		}
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

func (db *MemoryDB) GetMovieByID(id uint64) (entity.Movie, bool) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	m, ok := db.movies[id]
	return m, ok
}

func (db *MemoryDB) UpdateMovieByID(id uint64, updated entity.Movie) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	old, ok := db.movies[id]
	if !ok {
		return false
	}
	updated.ID = id
	if updated.Title != old.Title {
		delete(db.titleToID, old.Title)
		db.titleToID[updated.Title] = id
	}
	db.movies[id] = updated
	return true
}

func (db *MemoryDB) DeleteMovieByID(id uint64) bool {
	db.mu.Lock()
	defer db.mu.Unlock()
	old, ok := db.movies[id]
	if !ok {
		return false
	}
	delete(db.movies, id)
	delete(db.characters, id)
	if old.Title != "" {
		delete(db.titleToID, old.Title)
	}
	return true
}
