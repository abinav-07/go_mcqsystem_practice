package routes

import (
	"github/abinav-07/mcq-test/api/controllers"
	"github/abinav-07/mcq-test/api/middlewares"
	"github/abinav-07/mcq-test/router"
)

type AuthRoutes struct {
	router         router.Router
	middlewares    middlewares.AuthMW
	authController controllers.AuthController
}

// Group Auth router
func NewAuthRoutes(
	router router.Router,
	authController controllers.AuthController,
) AuthRoutes {
	return AuthRoutes{
		router:         router,
		authController: authController,
	}
}

// Setup Auth routes
func (i AuthRoutes) Setup() {
	auths := i.router.Gin.Group("/auth")

	//Grouped Auth routes
	auths.POST("login", i.authController.LoginUser)
	auths.POST("register", i.authController.RegisterUser)
}
