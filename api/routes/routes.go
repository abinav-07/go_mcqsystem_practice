package routes

import "go.uber.org/fx"

// Export Moduel
var Module = fx.Options(
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewUserRoutes),
	fx.Provide(NewRoutes),
	fx.Provide(NewTestRoutes),
)

// Routes contains multiple routes of interface type
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// Constructor Sets up all routes
func NewRoutes(userRoutes UserRoutes, authRoutes AuthRoutes, testRoutes TestRoutes) Routes {
	return Routes{
		userRoutes,
		authRoutes,
		testRoutes,
	}
}

// Setup all routes
func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
