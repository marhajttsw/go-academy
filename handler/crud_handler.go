package handler

import (
	"net/http"
	"project/db"
	"project/entity"
	"strconv"

	"go.uber.org/zap"
)

func RegisterCRUDHandlers(db *db.MemoryDB, mux *http.ServeMux, log *zap.Logger) {
	// movie endpoints
	mux.HandleFunc("/movies/add", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		id, _ := strconv.ParseUint(q.Get("id"), 10, 64)
		year, _ := strconv.Atoi(q.Get("year"))
		movie := entity.Movie{ID: id, Title: q.Get("title"), Year: year}
		db.AddMovie(movie)
		log.Info("AddMovie", zap.Uint64("id", movie.ID), zap.String("title", movie.Title), zap.Int("year", movie.Year))
	})

	mux.HandleFunc("/movies/get", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		m, ok := db.GetMovie(title)
		log.Info("GetMovie", zap.String("title", title), zap.Bool("found", ok), zap.Uint64("id", m.ID), zap.Int("year", m.Year))
	})

	mux.HandleFunc("/movies/update", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		title := q.Get("title")
		id, _ := strconv.ParseUint(q.Get("id"), 10, 64)
		year, _ := strconv.Atoi(q.Get("year"))
		updated := entity.Movie{ID: id, Title: title, Year: year}
		ok := db.UpdateMovie(title, updated)
		log.Info("UpdateMovie", zap.String("title", title), zap.Bool("ok", ok), zap.Uint64("id", updated.ID), zap.Int("year", updated.Year))
	})

	mux.HandleFunc("/movies/delete", func(w http.ResponseWriter, r *http.Request) {
		title := r.URL.Query().Get("title")
		ok := db.DeleteMovie(title)
		log.Info("DeleteMovie", zap.String("title", title), zap.Bool("ok", ok))
	})

	mux.HandleFunc("/movies/list", func(w http.ResponseWriter, r *http.Request) {
		movies := db.ListMovies()
		for _, m := range movies {
			log.Info("ListMovies", zap.Uint64("id", m.ID), zap.String("title", m.Title), zap.Int("year", m.Year))
		}
	})

	// character endpoints
	mux.HandleFunc("/characters/add", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		movie := q.Get("movie")
		name := q.Get("name")
		ok := db.AddCharacter(movie, entity.Character{Movie: movie, Name: name})
		log.Info("AddCharacter", zap.String("movie", movie), zap.String("name", name), zap.Bool("ok", ok))
	})

	mux.HandleFunc("/characters/get", func(w http.ResponseWriter, r *http.Request) {
		movie := r.URL.Query().Get("movie")
		chars := db.GetCharacters(movie)
		log.Info("GetCharacters", zap.String("movie", movie), zap.Int("count", len(chars)))
		for _, c := range chars {
			log.Info("Character", zap.String("movie", movie), zap.String("name", c.Name))
		}
	})

	mux.HandleFunc("/characters/update", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		movie := q.Get("movie")
		oldName := q.Get("old")
		newName := q.Get("new")
		ok := db.UpdateCharacter(movie, oldName, newName)
		log.Info("UpdateCharacter", zap.String("movie", movie), zap.String("old", oldName), zap.String("new", newName), zap.Bool("ok", ok))
	})

	mux.HandleFunc("/characters/delete", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		movie := q.Get("movie")
		name := q.Get("name")
		ok := db.DeleteCharacter(movie, name)
		log.Info("DeleteCharacter", zap.String("movie", movie), zap.String("name", name), zap.Bool("ok", ok))
	})
}
