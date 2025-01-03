// Package http defines functionality for default http server with minimum
// required configurations.
package http

import (
	"context"
	"fmt"
	"maker-checker/config"
	"maker-checker/pkg/db"

	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// Server represents the http server of the service.
type Server struct {
	GinRouter *gin.Engine
	DB        *db.DB
	Port      uint16
	cancel    context.CancelFunc
}

// New returns a new instance of http server.
func New(dbApp *db.DB, cfg config.AppConfig, cancel context.CancelFunc) *Server {
	// setup gin server.
	g := gin.Default()

	return &Server{
		GinRouter: g,
		DB:        dbApp,
		Port:      cfg.App.Port,
		cancel:    cancel,
	}
}

// Start starts the server.
func (server *Server) Start(ctx context.Context) error {
	errorGroup, ctx := errgroup.WithContext(ctx)
	appPort := server.Port

	errorGroup.Go(func() error {
		return server.startServer(ctx, appPort)
	})

	err := errorGroup.Wait()

	if err != nil {
		server.Stop()
	}

	return err
}

func (server *Server) startServer(ctx context.Context, port uint16) error {
	errCh := make(chan error)

	go func() {
		if err := server.GinRouter.Run(fmt.Sprintf(":%d", port)); err != nil {
			server.cancel()

			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():

	case err := <-errCh:
		return err
	}

	return nil
}

// Stop stops the server.
func (server *Server) Stop() {
	server.cancel()
}
