package echo_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	appapi "project/internal/api"
	"project/internal/db"
	appecho "project/internal/echo"
	"project/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func newTestServer() *echo.Echo {
	e := appecho.New(nil)
	database := db.New()
	// nil resty client is fine for tests
	h := handler.NewApiHandler(database, nil)
	appapi.RegisterHandlers(e, appapi.NewStrictHandler(h, nil))
	return e
}

func TestEchoRoutes_MoviesCRUD(t *testing.T) {
	e := newTestServer()

	// Create movie
	body, _ := json.Marshal(appapi.Movie{Name: "Interstellar", Year: 2014})
	req := httptest.NewRequest(http.MethodPost, "/movies", bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusCreated, rec.Code)

	var created appapi.Movie
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &created))
	require.Greater(t, created.Id, uint64(0))

	// Get list
	req = httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// Get by ID
	req = httptest.NewRequest(http.MethodGet, "/movies/"+jsonID(created.Id), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	// Delete
	req = httptest.NewRequest(http.MethodDelete, "/movies/"+jsonID(created.Id), nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusNoContent, rec.Code)
}

func jsonID(id uint64) string {
	return fmtUint64(id)
}

func fmtUint64(v uint64) string {
	if v == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
