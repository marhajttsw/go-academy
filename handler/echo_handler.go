package handler

import (
	"io"
	"net/http"

	"go.uber.org/zap"
)

// EchoHandler is an http.Handler that copies its request body
// back to the response.
type EchoHandler struct {
	log *zap.Logger
}

// NewEchoHandler builds a new EchoHandler.
func NewEchoHandler(log *zap.Logger) *EchoHandler {
	return &EchoHandler{log: log}
}

// ServeHTTP handles an HTTP request to the /echo endpoint.
func (h *EchoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	n, err := io.Copy(w, r.Body)
	if err != nil {
		h.log.Warn("Failed to handle request", zap.Error(err))
		return
	}
	h.log.Info("handled", zap.Int64("bytes", n))
}
