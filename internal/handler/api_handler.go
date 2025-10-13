package handler

import (
	"context"
	"project/internal/api"
	"project/internal/db"
	"project/internal/entity"
)

type ApiHandler struct {
	db *db.MemoryDB
}

func NewApiHandler(db *db.MemoryDB) *ApiHandler {
	return &ApiHandler{db: db}
}

// Movies

func (h *ApiHandler) GetMovies(ctx context.Context, _ api.GetMoviesRequestObject) (api.GetMoviesResponseObject, error) {
	movies := h.db.ListMovies()
	out := make([]api.Movie, 0, len(movies))
	for _, m := range movies {
		out = append(out, api.Movie{Id: m.ID, Name: m.Title, Year: m.Year})
	}
	return api.GetMovies200JSONResponse(out), nil
}

func (h *ApiHandler) PostMovies(ctx context.Context, req api.PostMoviesRequestObject) (api.PostMoviesResponseObject, error) {
	if req.Body == nil {
		return api.PostMovies400Response{}, nil
	}
	m := *req.Body
	ent := entity.Movie{ID: m.Id, Title: m.Name, Year: m.Year}
	h.db.AddMovie(ent)
	if ent.ID == 0 {
		if got, ok := h.db.GetMovie(ent.Title); ok {
			ent.ID = got.ID
		}
	}
	return api.PostMovies201JSONResponse(api.Movie{Id: ent.ID, Name: ent.Title, Year: ent.Year}), nil
}

func (h *ApiHandler) DeleteMoviesMovieId(ctx context.Context, req api.DeleteMoviesMovieIdRequestObject) (api.DeleteMoviesMovieIdResponseObject, error) {
	if h.db.DeleteMovieByID(req.MovieId) {
		return api.DeleteMoviesMovieId204Response{}, nil
	}
	return api.DeleteMoviesMovieId404Response{}, nil
}

func (h *ApiHandler) GetMoviesMovieId(ctx context.Context, req api.GetMoviesMovieIdRequestObject) (api.GetMoviesMovieIdResponseObject, error) {
	if m, ok := h.db.GetMovieByID(req.MovieId); ok {
		return api.GetMoviesMovieId200JSONResponse(api.Movie{Id: m.ID, Name: m.Title, Year: m.Year}), nil
	}
	return api.GetMoviesMovieId404Response{}, nil
}

func (h *ApiHandler) PutMoviesMovieId(ctx context.Context, req api.PutMoviesMovieIdRequestObject) (api.PutMoviesMovieIdResponseObject, error) {
	if req.Body == nil {
		return api.PutMoviesMovieId400Response{}, nil
	}
	m := *req.Body
	ent := entity.Movie{ID: req.MovieId, Title: m.Name, Year: m.Year}
	if ok := h.db.UpdateMovieByID(req.MovieId, ent); ok {
		return api.PutMoviesMovieId200JSONResponse(api.Movie{Id: ent.ID, Name: ent.Title, Year: ent.Year}), nil
	}
	return api.PutMoviesMovieId404Response{}, nil
}

// Characters

func (h *ApiHandler) GetCharacters(ctx context.Context, _ api.GetCharactersRequestObject) (api.GetCharactersResponseObject, error) {
	var out []api.Character
	for _, movie := range h.db.ListMovies() {
		chars := h.db.GetCharacters(movie.Title)
		for _, c := range chars {
			var idPtr *uint64
			if c.CharacterId != 0 {
				v := c.CharacterId
				idPtr = &v
			}
			out = append(out, api.Character{CharacterId: idPtr, Movie: c.Movie, MovieId: c.MovieID, Name: c.Name})
		}
	}
	return api.GetCharacters200JSONResponse(out), nil
}

func (h *ApiHandler) PostCharacters(ctx context.Context, req api.PostCharactersRequestObject) (api.PostCharactersResponseObject, error) {
	if req.Body == nil {
		return api.PostCharacters400Response{}, nil
	}
	c := *req.Body
	ent := entity.Character{Movie: c.Movie, Name: c.Name}
	if ok := h.db.AddCharacter(c.Movie, ent); !ok {
		return api.PostCharacters400Response{}, nil
	}
	chars := h.db.GetCharacters(c.Movie)
	var assigned *uint64
	if len(chars) > 0 {
		v := chars[len(chars)-1].CharacterId
		if v != 0 {
			assigned = &v
		}
	}
	return api.PostCharacters201JSONResponse(api.Character{CharacterId: assigned, Movie: c.Movie, MovieId: ent.MovieID, Name: c.Name}), nil
}

func (h *ApiHandler) DeleteCharactersCharacterId(ctx context.Context, req api.DeleteCharactersCharacterIdRequestObject) (api.DeleteCharactersCharacterIdResponseObject, error) {
	if h.db.DeleteCharacterByID(req.CharacterId) {
		return api.DeleteCharactersCharacterId204Response{}, nil
	}
	return api.DeleteCharactersCharacterId404Response{}, nil
}

func (h *ApiHandler) GetCharactersCharacterId(ctx context.Context, req api.GetCharactersCharacterIdRequestObject) (api.GetCharactersCharacterIdResponseObject, error) {
	if c, ok := h.db.GetCharacterByID(req.CharacterId); ok {
		var idPtr *uint64
		if c.CharacterId != 0 {
			v := c.CharacterId
			idPtr = &v
		}
		return api.GetCharactersCharacterId200JSONResponse(api.Character{CharacterId: idPtr, Movie: c.Movie, MovieId: c.MovieID, Name: c.Name}), nil
	}
	return api.GetCharactersCharacterId404Response{}, nil
}

func (h *ApiHandler) PutCharactersCharacterId(ctx context.Context, req api.PutCharactersCharacterIdRequestObject) (api.PutCharactersCharacterIdResponseObject, error) {
	if req.Body == nil {
		return api.PutCharactersCharacterId400Response{}, nil
	}
	b := *req.Body
	ent := entity.Character{Movie: b.Movie, Name: b.Name, MovieID: b.MovieId}
	if ok := h.db.UpdateCharacterByID(req.CharacterId, ent); ok {
		var idPtr *uint64
		v := req.CharacterId
		idPtr = &v
		return api.PutCharactersCharacterId200JSONResponse(api.Character{CharacterId: idPtr, Movie: ent.Movie, MovieId: ent.MovieID, Name: ent.Name}), nil
	}
	return api.PutCharactersCharacterId404Response{}, nil
}
