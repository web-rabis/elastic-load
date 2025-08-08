package cherver

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 30 * time.Second
	DefaultShutdownTimeout = 10 * time.Second
)

type Server struct {
	address           string
	certFile, keyFile string
	basePath          string

	readTimeout     time.Duration
	writeTimeout    time.Duration
	shutdownTimeout time.Duration

	plugins   []Plugin
	mws       chi.Middlewares
	resources []Resource

	server *http.Server
}

func New(opts ...Option) *Server {
	srv := &Server{}

	for i := range opts {
		opts[i](srv)
	}

	if srv.readTimeout <= 0 {
		srv.readTimeout = DefaultReadTimeout
	}

	if srv.writeTimeout <= 0 {
		srv.writeTimeout = DefaultWriteTimeout
	}

	if srv.shutdownTimeout <= 0 {
		srv.shutdownTimeout = DefaultShutdownTimeout
	}

	if srv.basePath == "" {
		srv.basePath = "/"
	}

	return srv
}

func (srv *Server) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Route(srv.basePath, func(r chi.Router) {
		r.Use(srv.mws...)

		for _, res := range srv.resources {
			r.Mount(res.Path(), res.Routes())
		}
	})

	return r
}

func (srv *Server) Run(ctx context.Context) (err error) {
	router := srv.setupRouter()

	srv.server = &http.Server{
		Addr:         srv.address,
		Handler:      router,
		ReadTimeout:  srv.readTimeout,
		WriteTimeout: srv.writeTimeout,
	}

	if len(srv.plugins) > 0 {
		log.Printf("[INFO] running cherver plugins...")
	}

	for i := range srv.plugins {
		if err := srv.plugins[i].Exec(ctx, router, srv.server); err != nil {
			return fmt.Errorf("plugin exec: %w", err)
		}
	}

	// prepare graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer func(f func() error) {
		gracefulErr := f()
		if err == nil {
			err = gracefulErr
		}
	}(srv.graceful(ctx))
	defer cancel()

	log.Printf("[INFO] serving HTTP on \"%s\"", srv.address)

	if srv.certFile == "" && srv.keyFile == "" {
		err = srv.server.ListenAndServe()
	} else {
		err = srv.server.ListenAndServeTLS(srv.certFile, srv.keyFile)
	}

	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (srv *Server) graceful(ctx context.Context) func() error {
	errChan := make(chan error, 1)

	go func() {
		defer close(errChan)
		<-ctx.Done()
		log.Printf("[INFO] shutting down HTTP server")

		ctx, cancel := context.WithTimeout(context.Background(), srv.shutdownTimeout)
		defer cancel()

		if err := srv.server.Shutdown(ctx); err != nil {
			errChan <- fmt.Errorf("HTTP server Shutdown: %v", err)
		}

		log.Println("[INFO] Cherver has processed all idle connections")
	}()

	return func() error {
		return <-errChan
	}
}
