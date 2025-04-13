package service_provider

import (
	"go.uber.org/fx"

	jwt_auth_repository "go-api-docker/internal/common/security/auth/infrastructure/jwt_auth_repository"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
)

var AuthServiceProvider = fx.Options(
	fx.Provide(
		fx.Annotate(
			jwt_auth_repository.NewJwtAuthRepository,
			fx.As(new(jwt_auth_repository.JwtAuthRepositoryInterface)),
		),
		fx.Annotate(
			jwt_auth.NewJwtAuthManager,
			fx.As(new(jwt_auth.JwtAuthManagerInterface)),
		),
	),
	fx.Invoke(
		jwt_auth_repository.NewJwtAuthRepository,
		jwt_auth.NewJwtAuthManager,
	),
)
