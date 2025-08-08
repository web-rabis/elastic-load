package http

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// FilesResource для раздачи статичных файлов
type FilesResource struct {
	path, filesDir string
}

func NewFilesResource(path, filesPath string) *FilesResource {
	return &FilesResource{
		path:     path,
		filesDir: filesPath,
	}
}

func (fr *FilesResource) Routes() chi.Router {
	r := chi.NewRouter()
	filesRoot := http.Dir(fr.filesDir)

	newFileServer(r, "/", filesRoot)

	return r
}

func (fr *FilesResource) Path() string {
	return fr.path
}

// NewFileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func newFileServer(r chi.Router, path string, root http.FileSystem) {
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
