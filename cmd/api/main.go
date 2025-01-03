package main

import (
	"context"
	"log"
	"maker-checker/config"
	"maker-checker/internal/middleware"
	"maker-checker/pkg/auth"
	"maker-checker/pkg/db"

	server "maker-checker/pkg/server/http"
)

func main() {
	// Load configuration
	cfg := config.Config()
	ctx, cancel := context.WithCancel(context.Background())

	dbApp, err := db.Init(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	defer dbApp.Stop()

	serverInstance := server.New(dbApp, *cfg, cancel)

	jwt := auth.New(cfg)
	middleware := middleware.New(jwt)
	api := serverInstance.GinRouter.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			messageRoutes(cfg, serverInstance.DB.DB(), v1, middleware)
			userRoutes(cfg, serverInstance.DB.DB(), v1)
		}
	}

	if err := serverInstance.Start(ctx); err != nil {
		panic(err)
	}
}
