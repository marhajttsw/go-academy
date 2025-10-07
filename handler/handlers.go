package handler

import "project/db"

func RegisterHandlers(db *db.MemoryDB) {
	RegisterMovieHandler(db)
	RegisterCharacterHandler(db)
}
