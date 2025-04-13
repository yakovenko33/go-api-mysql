package service_provider

import (
	"go.uber.org/fx"

	login_handler "go-api-docker/internal/go_crm/auth/application/service/login"
	refresh_tokens_handler "go-api-docker/internal/go_crm/auth/application/service/refresh_tokens"
	auth_controllers "go-api-docker/internal/go_crm/auth/ui/controllers"
	users_repository "go-api-docker/internal/go_crm/users/infrastructure/repositories/users_repository"
)

var AuthModuleServiceProvider = fx.Options(
	fx.Provide(
		fx.Annotate(
			users_repository.NewUsersRrepository,
			fx.As(new(users_repository.UsersRepositoryInterface)),
		),
		login_handler.NewLoginHandler,
		auth_controllers.NewAuthController,
		refresh_tokens_handler.NewLoginRefreshTokensHandler,
	),
	fx.Invoke(
		// users_repository.NewUsersRrepository,
		// auth_controllers.NewAuthController,
		// login_handler.NewLoginHandler,
		// refresh_tokens_handler.NewLoginRefreshTokensHandler,
		auth_controllers.RegisterAuthRoutes,
	),
)
