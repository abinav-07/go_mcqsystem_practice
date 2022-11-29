package routes

import (
	"github/abinav-07/mcq-test/api/controllers"
	"github/abinav-07/mcq-test/api/middlewares"
	"github/abinav-07/mcq-test/router"
)

type TestRoutes struct {
	router             router.Router
	testController     controllers.TestController
	questionController controllers.QuestionController
	middlewares        middlewares.AuthMW
	fbMiddleware       middlewares.FirebaseAuthMW
	trxMiddleware      middlewares.DBTransactionMW
}

// Group Test Routes
func NewTestRoutes(
	router router.Router,
	testController controllers.TestController,
	questionController controllers.QuestionController,
	middlewares middlewares.AuthMW,
	trxMiddleware middlewares.DBTransactionMW,
	fbMiddleware middlewares.FirebaseAuthMW,
) TestRoutes {
	return TestRoutes{
		router:             router,
		testController:     testController,
		questionController: questionController,
		middlewares:        middlewares,
		fbMiddleware:       fbMiddleware,
		trxMiddleware:      trxMiddleware,
	}
}

// Setup
func (i TestRoutes) Setup() {
	tests := i.router.Gin.Group("/test")
	tests.Use(i.middlewares.CheckJWT())
	tests.Use(i.middlewares.CheckAdmin())

	//Grouped Routes
	tests.POST("create", i.fbMiddleware.HandleAdmin(), i.trxMiddleware.HandleDBTransaction(), i.testController.CreateTests)
	tests.GET("/:testId", i.testController.GetTestDetails)
	tests.PATCH("/:testId/update", i.testController.UpdatePartial)
	tests.POST("/:testId/question/add", i.questionController.CreateQuestionAndAnswers)

}
