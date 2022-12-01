package routes

import (
	"github/abinav-07/mcq-test/api/controllers"
	"github/abinav-07/mcq-test/api/middlewares"
	"github/abinav-07/mcq-test/router"
)

type UserRoutes struct {
	router          router.Router
	userController  controllers.UserController
	testController  controllers.TestController
	adminController controllers.AdminController
	middlewares     middlewares.AuthMW
	trxMiddleware   middlewares.DBTransactionMW
	fbMiddleware    middlewares.FirebaseAuthMW
}

// Group user routes
func NewUserRoutes(
	router router.Router,
	userController controllers.UserController,
	testController controllers.TestController,
	adminController controllers.AdminController,
	middlewares middlewares.AuthMW,
	trxMiddleware middlewares.DBTransactionMW,
	fbMiddleware middlewares.FirebaseAuthMW,

) UserRoutes {
	return UserRoutes{
		router:          router,
		userController:  userController,
		testController:  testController,
		middlewares:     middlewares,
		trxMiddleware:   trxMiddleware,
		fbMiddleware:    fbMiddleware,
		adminController: adminController,
	}
}

// Setup user routes
func (i UserRoutes) Setup() {

	//Update User Details By admin
	i.router.Gin.PATCH("/admin/user/:userId",
		i.middlewares.CheckJWT(),
		i.middlewares.CheckAdmin(),
		i.fbMiddleware.HandleAdmin(),
		i.trxMiddleware.HandleDBTransaction(),
		i.adminController.UpdateUser)

	users := i.router.Gin.Group("/user")
	users.Use(i.middlewares.CheckJWT())
	//Grouped Routes
	users.DELETE("/:userId", i.trxMiddleware.HandleDBTransaction(), i.userController.DeleteUserById)

	users.GET("tests", i.userController.GetAvailableTests)
	users.GET("tests/:testId", i.testController.GetTestDetails)
	users.POST("tests/submit-test", i.userController.CreateTestReport)
}
