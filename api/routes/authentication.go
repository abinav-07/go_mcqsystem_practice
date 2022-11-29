package routes

import (
	"github/abinav-07/mcq-test/api/controllers"
	"github/abinav-07/mcq-test/api/middlewares"
	"github/abinav-07/mcq-test/router"
)

type AuthRoutes struct {
	router         router.Router
	middlewares    middlewares.DBTransactionMW
	authController controllers.AuthController
}

// Group Auth router
func NewAuthRoutes(
	router router.Router,
	authController controllers.AuthController,
	middlewares middlewares.DBTransactionMW,
) AuthRoutes {
	return AuthRoutes{
		router:         router,
		authController: authController,
		middlewares:    middlewares,
	}
}

// Setup Auth routes
func (i AuthRoutes) Setup() {
	auths := i.router.Gin.Group("/auth")

	//Grouped Auth routes
	auths.POST("login", i.middlewares.HandleDBTransaction(), i.authController.LoginUser)
	auths.POST("register", i.middlewares.HandleDBTransaction(), i.authController.RegisterUser)
}
