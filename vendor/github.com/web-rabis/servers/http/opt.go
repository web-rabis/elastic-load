package cherver

import (
	"net/http"
	"time"
)

type Option func(*Server)

// WithListenAddress sets listen address. By default its :80
func WithListenAddress(addr string) Option {
	return func(s *Server) {
		s.address = addr
	}
}

// WithCert sets server to HTTPS
func WithCert(certFile, keyFile string) Option {
	return func(s *Server) {
		s.certFile = certFile
		s.keyFile = keyFile
	}
}

// WithBasePath sets base path for all server. By default is "/"
func WithBasePath(basePath string) Option {
	return func(s *Server) {
		s.basePath = basePath
	}
}

// WithMiddlewares sets middlewares to all routes
func WithMiddlewares(mws ...func(http.Handler) http.Handler) Option {
	return func(s *Server) {
		s.mws = mws
	}
}

// WithResources sets resources
func WithResources(r ...Resource) Option {
	return func(s *Server) {
		s.resources = r
	}
}

func WithReadTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.readTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.writeTimeout = t
	}
}

func WithShutdownTimeout(t time.Duration) Option {
	return func(s *Server) {
		s.shutdownTimeout = t
	}
}

func WithPlugins(ps ...Plugin) Option {
	return func(s *Server) {
		s.plugins = ps
	}
}
