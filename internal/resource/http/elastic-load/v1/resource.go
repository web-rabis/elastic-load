package v1

import (
	"github.com/go-chi/chi/v5"
	"github.com/web-rabis/elastic-load/internal/manager/elk"
)

type Resource struct {
	path   string
	elkMan elk.IManager
}

func NewElasticLoadResource(basePath string, elkMan elk.IManager) *Resource {
	return &Resource{
		path:   basePath,
		elkMan: elkMan,
	}
}

func (res *Resource) Path() string {
	return res.path
}

func (res *Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Get("/full", res.full)
		r.Get("/full/stop", res.fullStop)
		r.Get("/full/status", res.fullStatus)
	})

	return r
}
