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

	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
)

var (
	once   sync.Once
	router *gin.Engine
)

func NewRouter(jwt_auth jwt_auth.JwtAuthManagerInterface) *gin.Engine {
	once.Do(func() {
		router = gin.Default()
		/*config := cors.DefaultConfig()
		config.AllowAllOrigins = true


		config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}
		config.AllowHeaders = []string{"Origin", "Content-Type"}
		router.Use(cors.New(config))*/

		router.Use(cors.New(cors.Config{
			AllowOrigins:     []string{"https://your-frontend-domain.com"},        // Разрешенные домены (домен вашего SPA)
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Разрешенные HTTP методы
			AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // Разрешенные заголовки
			AllowCredentials: true,                                                // Разрешить отправку cookie
		}))
	})

	return router
}

func StartServer(lc fx.Lifecycle, router *gin.Engine, logger *zap.Logger) {
	server := &http.Server{
		Addr:         ":8080",
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
