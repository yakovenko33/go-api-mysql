package service_provider

import (
	"go.uber.org/fx"

	auth_controllers "go-api-docker/internal/go_crm/auth/ui/controllers"
)

var AuthModuleServiceProvider = fx.Options(
	fx.Invoke(auth_controllers.RegisterAuthRoutes),
)
