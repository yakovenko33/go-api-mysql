package server

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	router *gin.Engine
)

func NewRouter() *gin.Engine {
	once.Do(func() {
		router = gin.Default()
		config := cors.DefaultConfig()
		config.AllowAllOrigins = true
		config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}
		config.AllowHeaders = []string{"Origin", "Content-Type"}
		router.Use(cors.New(config))
	})

	return router
}

func StartServer(lc fx.Lifecycle, router *gin.Engine, logger *zap.Logger) {
	server := &http.Server{
		Addr:         ":3000",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					logger.Error(fmt.Sprintf("Error starting server: %s\n", err))
				}
			}()
			fmt.Println("Gin server started on :8080")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Stopping Gin server...")
			return server.Shutdown(ctx)
		},
	})
}
