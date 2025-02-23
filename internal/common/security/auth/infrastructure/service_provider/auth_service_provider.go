package service_provider

import (
	"go.uber.org/fx"

	jwt_auth_repository "go-api-docker/internal/common/security/auth/infrastructure/jwt_auth_repository"
	jwt_auth "go-api-docker/internal/common/security/auth/jwt_auth"
)

var AuthServiceProvider = fx.Options(
	fx.Provide(
		jwt_auth_repository.NewJwtAuthRepository,
		fx.Annotate(
			func(r *jwt_auth_repository.JwtAuthRepository) *jwt_auth_repository.JwtAuthRepository {
				return r
			},
			fx.As(new(jwt_auth.JwtAuthManager)),
		),
	),
	fx.Invoke(jwt_auth_repository.NewJwtAuthRepository),
)
