package handler_test

import (
    "context"
    "testing"

    "project/internal/api"
    "project/internal/db"
    "project/internal/handler"

    "github.com/stretchr/testify/require"
)

func TestApiHandler_MoviesFlow(t *testing.T) {
    t.Parallel()
    database := db.New()
    h := handler.NewApiHandler(database, nil)
    ctx := context.Background()

    // List empty
    resp, err := h.GetMovies(ctx, api.GetMoviesRequestObject{})
    require.NoError(t, err)
    if r, ok := resp.(api.GetMovies200JSONResponse); ok {
        require.Len(t, r, 0)
    } else {
        t.Fatalf("unexpected response type: %T", resp)
    }

    // Create movie
    createResp, err := h.PostMovies(ctx, api.PostMoviesRequestObject{Body: &api.Movie{Name: "Inception", Year: 2010}})
    require.NoError(t, err)
    var movieID uint64
    if r, ok := createResp.(api.PostMovies201JSONResponse); ok {
        require.Greater(t, r.Id, uint64(0))
        require.Equal(t, "Inception", r.Name)
        movieID = r.Id
    } else {
        t.Fatalf("unexpected response type: %T", createResp)
    }

    // List has one
    resp, err = h.GetMovies(ctx, api.GetMoviesRequestObject{})
    require.NoError(t, err)
    if r, ok := resp.(api.GetMovies200JSONResponse); ok {
        require.Len(t, r, 1)
        require.Equal(t, movieID, r[0].Id)
    } else {
        t.Fatalf("unexpected response type: %T", resp)
    }

    // Get by ID
    getResp, err := h.GetMoviesMovieId(ctx, api.GetMoviesMovieIdRequestObject{MovieId: movieID})
    require.NoError(t, err)
    if r, ok := getResp.(api.GetMoviesMovieId200JSONResponse); ok {
        require.Equal(t, movieID, r.Id)
        require.Equal(t, "Inception", r.Name)
    } else {
        t.Fatalf("unexpected response type: %T", getResp)
    }

    // Update by ID
    updResp, err := h.PutMoviesMovieId(ctx, api.PutMoviesMovieIdRequestObject{MovieId: movieID, Body: &api.Movie{Name: "Inception (Upd)", Year: 2011}})
    require.NoError(t, err)
    if r, ok := updResp.(api.PutMoviesMovieId200JSONResponse); ok {
        require.Equal(t, movieID, r.Id)
        require.Equal(t, "Inception (Upd)", r.Name)
        require.Equal(t, 2011, r.Year)
    } else {
        t.Fatalf("unexpected response type: %T", updResp)
    }

    // Delete by ID
    delResp, err := h.DeleteMoviesMovieId(ctx, api.DeleteMoviesMovieIdRequestObject{MovieId: movieID})
    require.NoError(t, err)
    if _, ok := delResp.(api.DeleteMoviesMovieId204Response); !ok {
        t.Fatalf("expected 204 response, got: %T", delResp)
    }

    // Get non-existing
    notFound, err := h.GetMoviesMovieId(ctx, api.GetMoviesMovieIdRequestObject{MovieId: movieID})
    require.NoError(t, err)
    if _, ok := notFound.(api.GetMoviesMovieId404Response); !ok {
        t.Fatalf("expected 404 response, got: %T", notFound)
    }
}

func TestApiHandler_CharactersFlow(t *testing.T) {
    t.Parallel()
    database := db.New()
    h := handler.NewApiHandler(database, nil)
    ctx := context.Background()

    // Seed a movie
    created, err := h.PostMovies(ctx, api.PostMoviesRequestObject{Body: &api.Movie{Name: "Inception", Year: 2010}})
    require.NoError(t, err)
    var movieID uint64
    if r, ok := created.(api.PostMovies201JSONResponse); ok {
        movieID = r.Id
    } else {
        t.Fatalf("unexpected response type: %T", created)
    }

    // Create character (avoid Star Wars to skip external SWAPI)
    cCreate, err := h.PostCharacters(ctx, api.PostCharactersRequestObject{Body: &api.Character{Movie: "Inception", Name: "Dom Cobb"}})
    require.NoError(t, err)
    var charID uint64
    if r, ok := cCreate.(api.PostCharacters201JSONResponse); ok {
        require.Equal(t, "Inception", r.Movie)
        require.Equal(t, movieID, r.MovieId)
        require.NotNil(t, r.CharacterId)
        charID = *r.CharacterId
        require.Greater(t, charID, uint64(0))
    } else {
        t.Fatalf("unexpected response type: %T", cCreate)
    }

    // List characters
    listResp, err := h.GetCharacters(ctx, api.GetCharactersRequestObject{})
    require.NoError(t, err)
    if r, ok := listResp.(api.GetCharacters200JSONResponse); ok {
        require.Len(t, r, 1)
        require.Equal(t, charID, *r[0].CharacterId)
    } else {
        t.Fatalf("unexpected response type: %T", listResp)
    }

    // Get character by ID
    getResp, err := h.GetCharactersCharacterId(ctx, api.GetCharactersCharacterIdRequestObject{CharacterId: charID})
    require.NoError(t, err)
    if r, ok := getResp.(api.GetCharactersCharacterId200JSONResponse); ok {
        require.NotNil(t, r.CharacterId)
        require.Equal(t, charID, *r.CharacterId)
        require.Equal(t, "Dom Cobb", r.Name)
    } else {
        t.Fatalf("unexpected response type: %T", getResp)
    }

    // Update character by ID
    updResp, err := h.PutCharactersCharacterId(ctx, api.PutCharactersCharacterIdRequestObject{CharacterId: charID, Body: &api.Character{Movie: "Inception", Name: "Arthur", MovieId: movieID}})
    require.NoError(t, err)
    if r, ok := updResp.(api.PutCharactersCharacterId200JSONResponse); ok {
        require.NotNil(t, r.CharacterId)
        require.Equal(t, charID, *r.CharacterId)
        require.Equal(t, "Arthur", r.Name)
    } else {
        t.Fatalf("unexpected response type: %T", updResp)
    }

    // Delete character by ID
    delResp, err := h.DeleteCharactersCharacterId(ctx, api.DeleteCharactersCharacterIdRequestObject{CharacterId: charID})
    require.NoError(t, err)
    if _, ok := delResp.(api.DeleteCharactersCharacterId204Response); !ok {
        t.Fatalf("expected 204 response, got: %T", delResp)
    }

    // Get non-existing
    notFound, err := h.GetCharactersCharacterId(ctx, api.GetCharactersCharacterIdRequestObject{CharacterId: charID})
    require.NoError(t, err)
    if _, ok := notFound.(api.GetCharactersCharacterId404Response); !ok {
        t.Fatalf("expected 404 response, got: %T", notFound)
    }
}

func TestApiHandler_NilBodies(t *testing.T) {
    t.Parallel()
    database := db.New()
    h := handler.NewApiHandler(database, nil)
    ctx := context.Background()

    // Movies
    r1, err := h.PostMovies(ctx, api.PostMoviesRequestObject{Body: nil})
    require.NoError(t, err)
    if _, ok := r1.(api.PostMovies400Response); !ok {
        t.Fatalf("expected 400 response, got: %T", r1)
    }
    r2, err := h.PutMoviesMovieId(ctx, api.PutMoviesMovieIdRequestObject{MovieId: 1, Body: nil})
    require.NoError(t, err)
    if _, ok := r2.(api.PutMoviesMovieId400Response); !ok {
        t.Fatalf("expected 400 response, got: %T", r2)
    }

    // Characters
    r3, err := h.PostCharacters(ctx, api.PostCharactersRequestObject{Body: nil})
    require.NoError(t, err)
    if _, ok := r3.(api.PostCharacters400Response); !ok {
        t.Fatalf("expected 400 response, got: %T", r3)
    }
    r4, err := h.PutCharactersCharacterId(ctx, api.PutCharactersCharacterIdRequestObject{CharacterId: 1, Body: nil})
    require.NoError(t, err)
    if _, ok := r4.(api.PutCharactersCharacterId400Response); !ok {
        t.Fatalf("expected 400 response, got: %T", r4)
    }
}
