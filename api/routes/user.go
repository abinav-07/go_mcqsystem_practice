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
	trxMiddleware  middlewares.DBTransactionMW
}

// Group user routes
func NewUserRoutes(
	router router.Router,
	userController controllers.UserController,
	testController controllers.TestController,
	middlewares middlewares.AuthMW,
	trxMiddleware middlewares.DBTransactionMW,

) UserRoutes {
	return UserRoutes{
		router:         router,
		userController: userController,
		testController: testController,
		middlewares:    middlewares,
		trxMiddleware:  trxMiddleware,
	}
}

// Setup user routes
func (i UserRoutes) Setup() {
	users := i.router.Gin.Group("/user")
	users.Use(i.middlewares.CheckJWT())
	//Grouped Routes
	users.DELETE("/:userId", i.trxMiddleware.HandleDBTransaction(), i.userController.DeleteUserById)

	users.GET("tests", i.userController.GetAvailableTests)
	users.GET("tests/:testId", i.testController.GetTestDetails)
	users.POST("tests/submit-test", i.userController.CreateTestReport)
}
