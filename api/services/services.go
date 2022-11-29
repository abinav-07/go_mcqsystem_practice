package services

import "go.uber.org/fx"

// Export
var Module = fx.Options(
	fx.Provide(NewJWTAuthService),
	fx.Provide(NewUserService),
	fx.Provide(NewRoleService),
	fx.Provide(NewTestService),
	fx.Provide(NewQuestionService),
	fx.Provide(NewUserTestService),
	fx.Provide(NewFirebaseService),
	fx.Provide(NewFireStoreClient),
)
