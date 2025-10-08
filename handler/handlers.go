package handler

import (
	"net/http"
	"project/db"
	"go.uber.org/zap"
)

func RegisterHandlers(db *db.MemoryDB, mux *http.ServeMux, echo *EchoHandler, log *zap.Logger) {
	RegisterMovieHandler(db)      // unchanged
	RegisterCharacterHandler(db)  // unchanged
	RegisterCRUDHandlers(db, mux, log)
	mux.Handle("/echo", echo)
}
