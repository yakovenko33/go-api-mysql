package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	database "go-api-docker/internal/common/database"
	logging "go-api-docker/internal/common/logging"
	access_control_model "go-api-docker/internal/common/security/access_control_models"
	users_controllers "go-api-docker/internal/go_crm/users/ui/controllers"

	"go.uber.org/zap"
)

func NewRouter() *gin.Engine {
	router := gin.Default()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	router.Use(cors.New(config))
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

func coreApp() *fx.App {
	app := fx.New(
		fx.Provide(
			logging.InitLogging,
			database.ProvideDBConnection,
			access_control_model.InitAccessControlModel,
			NewRouter,
		),
		fx.Invoke(
			users_controllers.RegisterUserRoutes,
			StartServer,
		),
	)
	return app
}

func main() {
	app := coreApp()
	app.Run()
}
