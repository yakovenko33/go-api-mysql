package main

import (
	"go.uber.org/fx"

	database "go-api-docker/internal/common/database"
	logging "go-api-docker/internal/common/logging"
	access_control_model "go-api-docker/internal/common/security/access_control_models"
	auth_service_provider "go-api-docker/internal/common/security/auth/infrastructure/service_provider"
	server "go-api-docker/internal/common/server"
	auth_module_service_provider "go-api-docker/internal/go_crm/auth/infrastructure/service_provider"
	users_service_provider "go-api-docker/internal/go_crm/users/infrastructure/service_provider"
)

func coreApp() *fx.App {
	app := fx.New(
		fx.Provide(
			logging.InitLogging,
			database.ProvideDBConnection,
			access_control_model.InitAccessControlModel,
			server.NewRouter,
		),
		auth_service_provider.AuthServiceProvider,
		users_service_provider.UsersServiceProvider,
		auth_module_service_provider.AuthModuleServiceProvider,
		fx.Invoke(
			server.StartServer,
		),
	)
	return app
}

func main() {
	app := coreApp()
	app.Run()
}
