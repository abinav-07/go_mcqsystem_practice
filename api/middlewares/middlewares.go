package middlewares

import "go.uber.org/fx"

// Export Middlewares
var Module = fx.Options(
	fx.Provide(NewAuthMiddlware),
	fx.Provide(NewDBTransactionMiddleware),
	fx.Provide(NewFirebaseAuthMiddleware),
)
