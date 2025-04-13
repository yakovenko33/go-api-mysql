package service_provider

import (
	"go.uber.org/fx"

	users_controllers "go-api-docker/internal/go_crm/users/ui/controllers"
)

var UsersServiceProvider = fx.Options(
	fx.Invoke(
		users_controllers.RegisterUserRoutes,
	),
)
