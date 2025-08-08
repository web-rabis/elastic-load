package cherver

import "github.com/go-chi/chi/v5"

type Resource interface {
	Routes() chi.Router
	// Path - base path of your resource
	Path() string
}
