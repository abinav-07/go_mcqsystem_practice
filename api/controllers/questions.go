package controllers

import (
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	testService     services.TestService
	questionService services.QuestionService
}

// Constructor
func NewQuestionController(
	testService services.TestService,
	questionService services.QuestionService,
) QuestionController {
	return QuestionController{
		testService:     testService,
		questionService: questionService,
	}
}

func (qc QuestionController) CreateQuestionAndAnswers(ctx *gin.Context) {
	//Empty Struct
	reqBody := struct{ models.Question }{}

	testIdParam := ctx.Param("testId")

	if testIdParam == "" {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": true, "message": "Invalid testId on search filter!"})
		ctx.Abort()
		return
	}

	testId, _ := strconv.ParseUint(testIdParam, 10, 32)
	testIdParamUint := uint(testId)

	//Bind Body to test struct
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	reqBody.Question.TestID = testIdParamUint

	//Check if it is a valid test
	_, testErr := qc.testService.GetById(testIdParamUint)
	if testErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, " message": testErr})
		ctx.Abort()
		return
	}

	//Create Question and Options for the test
	createQuestions, createErr := qc.questionService.CreateQuestion(reqBody.Question)
	if createErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": createErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Questions and answers Created!", "data": createQuestions})

}
