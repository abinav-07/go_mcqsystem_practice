package controllers

import "go.uber.org/fx"

// Module exported for initializing application
var Module = fx.Options(
	fx.Provide(NewAuthController),
	fx.Provide(NewUserController),
	fx.Provide(NewTestController),
	fx.Provide(NewQuestionController),
	fx.Provide(NewAdminController),
)
