package cherver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Plugin interface {
	Exec(ctx context.Context, router chi.Router, server *http.Server) error
}
