package routes

import (
	"github/abinav-07/mcq-test/api/controllers"
	"github/abinav-07/mcq-test/api/middlewares"
	"github/abinav-07/mcq-test/router"
)

type UserRoutes struct {
	router         router.Router
	userController controllers.UserController
	testController controllers.TestController
	middlewares    middlewares.AuthMW
}

// Group user routes
func NewUserRoutes(
	router router.Router,
	userController controllers.UserController,
	testController controllers.TestController,
	middlewares middlewares.AuthMW,

) UserRoutes {
	return UserRoutes{
		router:         router,
		userController: userController,
		testController: testController,
		middlewares:    middlewares,
	}
}

// Setup user routes
func (i UserRoutes) Setup() {
	users := i.router.Gin.Group("/user")
	users.Use(i.middlewares.CheckJWT())
	//Grouped Routes
	users.GET("tests", i.userController.GetAvailableTests)
	users.GET("tests/:testId", i.testController.GetTestDetails)
	users.POST("submit-test", i.userController.CreateTestReport)
}
