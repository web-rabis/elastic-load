package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const APIVersion = "v1"

// VersionResponse - ответ на запрос версии.
type VersionResponse struct {
	API     string `json:"api"`
	Version string `json:"version"`
}

// VersionResource - структура содержащая версию API и приложения.
type VersionResource struct {
	path, version string
}

func NewVersionResource(path, version string) *VersionResource {
	return &VersionResource{path: path, version: version}
}

func (vr VersionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.get)

	return r
}

func (vr VersionResource) Path() string {
	return vr.path
}

func (vr VersionResource) get(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, VersionResponse{
		API:     APIVersion,
		Version: vr.version,
	})
}
