package v1

import (
	"path/filepath"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Elastic-load API
// @version 1.0
// @description Pay System Config API
// @BasePath /elastic-load/api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
type SwaggerResource struct {
	path      string
	basePath  string
	filesPath string
}

func NewSwaggerResource(path, basePath, filesPath string) SwaggerResource {
	return SwaggerResource{
		path:      path,
		basePath:  basePath,
		filesPath: filesPath,
	}
}

func (sr SwaggerResource) Path() string {
	return sr.path
}

func (sr SwaggerResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/*", httpSwagger.Handler(
		httpSwagger.URL(filepath.Join(sr.basePath, sr.filesPath, "swagger.json")),
	))
	return r
}
