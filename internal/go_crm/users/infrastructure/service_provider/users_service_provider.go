package service_provider

import (
	"go.uber.org/fx"

	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	auth_controllers "go-api-docker/internal/go_crm/auth/ui/controllers"
	users_repository "go-api-docker/internal/go_crm/users/infrastructure/repositories/users_repository"
	users_controllers "go-api-docker/internal/go_crm/users/ui/controllers"
)

var UsersServiceProvider = fx.Options(
	fx.Provide(
		fx.Annotate(
			users_repository.NewUsersRrepository,
			fx.As(new(users_repository.UsersRepositoryInterface)),
		),
		login_handler.NewLoginHandler,
		auth_controllers.NewAuthController,
	),
	fx.Invoke(
		users_repository.NewUsersRrepository,
		auth_controllers.NewAuthController,
		users_controllers.RegisterUserRoutes,
		login_handler.NewLoginHandler,
	),
)
