package infrastructure

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewEnv),
	fx.Provide(NewMCQTestDatabase),
	fx.Provide(NewFireBaseApp),
	fx.Provide(NewFirebaseAuth),
	fx.Provide(NewFirestoreClient),
)
