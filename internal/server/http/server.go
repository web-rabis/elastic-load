package http

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/web-rabis/elastic-load/internal/adapter/database/drivers"
	"github.com/web-rabis/elastic-load/internal/manager/block"
	"github.com/web-rabis/elastic-load/internal/manager/catalog"
	"github.com/web-rabis/elastic-load/internal/manager/dictionary"
	"github.com/web-rabis/elastic-load/internal/manager/ebook"

	"github.com/web-rabis/elastic-load/internal/config"
	"github.com/web-rabis/elastic-load/internal/manager/elk"
	"github.com/web-rabis/elastic-load/internal/resource/http"
	elasticV1 "github.com/web-rabis/elastic-load/internal/resource/http/elastic-load/v1"
	swaggerV1 "github.com/web-rabis/elastic-load/internal/resource/http/swagger/v1"
	cherver "github.com/web-rabis/servers/http"
)

const (
	maxAge        = 300
	compressLevel = 5
)

func Run(serversCtx context.Context, opts *config.APIServer, ds drivers.DataStore, version string) error {

	blockMan, err := block.NewBlockManager(serversCtx, ds.Block())
	if err != nil {
		return err
	}
	dictMan, err := dictionary.NewDictionaryManager(serversCtx, ds.Dictionary())
	if err != nil {
		return err
	}
	catMan, err := catalog.NewCatalogManager(serversCtx, ds.Catalog())
	if err != nil {
		return err
	}
	ebookMan := ebook.NewEbookManager(ds.Ebook(), blockMan, dictMan, catMan)
	elkMan, err := elk.NewElkManager(opts, ebookMan)
	if err != nil {
		return err
	}

	resources := []cherver.Resource{
		http.NewVersionResource("/version", version),
		http.NewFilesResource("/files", opts.ServerConfig.FilesDir),
		swaggerV1.NewSwaggerResource("/swagger", opts.ServerConfig.BasePath, "/files"),
		elasticV1.NewElasticLoadResource("/api/v1/elastic-load", elkMan),
	}
	httpSrv := cherver.New(
		cherver.WithListenAddress(opts.ServerConfig.ListenAddr),
		cherver.WithCert(opts.ServerConfig.CertFile, opts.ServerConfig.KeyFile),
		cherver.WithResources(resources...),
		cherver.WithMiddlewares(middlewares(opts)...))

	return httpSrv.Run(serversCtx)
}

func middlewaresWithoutLogs(opts *config.APIServer) chi.Middlewares {
	return chi.Middlewares{
		middleware.NoCache,   // no-cache
		middleware.Recoverer, // управляемо обрабатывает паники и выдает stack trace при их возникновении
		middleware.RealIP,    // устанавливает RemoteAddr для каждого запроса с заголовками X-Forwarded-For или X-Real-IP
		middleware.NewCompressor(compressLevel).Handler,

		cors.Handler(cors.Options{
			AllowedOrigins:   allowedOrigins(opts.IsTesting),
			AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           maxAge, // Maximum value not ignored by any of major browsers
		})}
}

func middlewares(opts *config.APIServer) chi.Middlewares {
	return append(middlewaresWithoutLogs(opts), middleware.Logger)
}
func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}
