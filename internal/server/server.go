package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"
)

var (
	isHealthy int32
)

type Server struct {
	addr   string
	router *http.ServeMux
	logger *zap.Logger
	quit   chan os.Signal
}

func New(addr string, logger *zap.Logger) *Server {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	return &Server{
		addr:   addr,
		router: http.NewServeMux(),
		logger: logger,
		quit:   quit,
	}
}

func (s *Server) AddHandler(route string, handler http.Handler) {
	s.router.Handle(route, handler)
}

func (s *Server) Health(route string) {
	s.router.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("call for health check")

		if atomic.LoadInt32(&isHealthy) == 1 {
			w.WriteHeader(http.StatusOK)
			return
		}

		w.WriteHeader(http.StatusServiceUnavailable)
	}))
}

func (s *Server) ListenAndServe() {
	server := &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}

	done := make(chan bool)

	go func() {
		<-s.quit
		s.logger.Info("server is shutting down ...")
		atomic.StoreInt32(&isHealthy, 0)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(ctx); err != nil {
			s.logger.Fatal("could not gracefully shutdown the server", zap.Error(err))
		}

		close(done)
	}()

	atomic.StoreInt32(&isHealthy, 1)
	s.logger.Debug("start server ...", zap.String("addr", server.Addr))
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logger.Fatal("could not listen", zap.Error(err))
	}

	<-done
	s.logger.Debug("server stopped")
}

func (s *Server) Shutdown() {
	signal.Notify(s.quit, syscall.SIGINT)
}
