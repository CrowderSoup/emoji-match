package server

import (
	"context"
	"embed"
	"io/fs"
	"net/http"
	"time"

	"github.com/CrowderSoup/emoji-match/config"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type contextKey string

type Server struct {
	Address                 string
	Logger                  *zap.SugaredLogger
	mux                     *http.ServeMux
	http                    *http.Server
	readStoredProcedureFunc func(filePath string) (string, error)
	StaticAssets            embed.FS
}

type Middleware func(http.Handler) http.Handler

// NewServer returns a new server
func NewServer(config *config.Config, logger *zap.SugaredLogger, mux *http.ServeMux, staticAssets embed.FS) *Server {
	server := &Server{
		Address:      config.Address,
		Logger:       logger,
		mux:          mux,
		StaticAssets: staticAssets,
	}

	server.initRoutes()

	s := &http.Server{
		Addr:           config.Address,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	server.http = s

	return server
}

func (s *Server) initRoutes() {
	// Register static file path
	fs, err := fs.Sub(s.StaticAssets, "dist")
	if err != nil {
		s.Logger.Fatal(err)
	}
	staticHandler := http.FileServer(http.FS(fs))

	s.wireRoute("/healthcheck", http.HandlerFunc(s.healthcheck))
	s.wireRoute("/", staticHandler)
}

func (s *Server) errorHandler(w http.ResponseWriter, msg string, statusCode int) {
	s.Logger.Error(msg)

	w.WriteHeader(statusCode)
	w.Write([]byte(msg))
}

func (s *Server) sendResponse(w http.ResponseWriter, body string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Write([]byte(body))
}

func (s *Server) wireRoute(pattern string, handler http.Handler, middlewares ...Middleware) {
	// Inject some middleware we want everywhere
	telemetryMiddleware := otelhttp.NewMiddleware(pattern)
	loggingMiddleware := LoggingMiddleware(s.Logger)
	s.mux.Handle(
		pattern,
		s.useMiddlewares(
			handler,
			append(
				middlewares,
				telemetryMiddleware,
				CORSMiddleware,
				loggingMiddleware,
			)...,
		),
	)
}

func (s *Server) useMiddlewares(h http.Handler, middlewares ...Middleware) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- { // wrap inside-out
		h = middlewares[i](h)
	}
	return h
}

// Run starts our HTTP server
func Run(lc fx.Lifecycle, s *Server) {
	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				s.Logger.Infow("starting http server",
					"Address", s.Address,
				)

				go s.http.ListenAndServe()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.http.Shutdown(ctx)
				return nil
			},
		},
	)
}

// Module provided to fx
var Module = fx.Options(
	fx.Provide(
		http.NewServeMux,
		NewServer,
	),
	fx.Invoke(Run),
)
