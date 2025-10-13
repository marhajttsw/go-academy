package handler

import (
	"project/internal/db"

	"go.uber.org/zap"
)

func RegisterHandlers(db *db.MemoryDB, log *zap.Logger) {
	// RegisterMovieHandler(db)     // unchanged
	// RegisterCharacterHandler(db) // unchanged
	//RegisterCRUDHandlers(db, mux, log)
	//mux.Handle("/echo", echo)
}
