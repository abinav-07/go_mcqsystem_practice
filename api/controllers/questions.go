package controllers

import (
	"encoding/json"
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionController struct {
	testService      services.TestService
	questionService  services.QuestionService
	fireStoreService services.FireStoreService
}

// Constructor
func NewQuestionController(
	testService services.TestService,
	questionService services.QuestionService,
	fireStoreService services.FireStoreService,
) QuestionController {
	return QuestionController{
		testService:      testService,
		questionService:  questionService,
		fireStoreService: fireStoreService,
	}
}

func (qc QuestionController) CreateQuestionAndAnswers(ctx *gin.Context) {
	//Empty Struct
	reqBody := struct{ models.Question }{}

	testIdParam := ctx.Param("testId")

	if testIdParam == "" {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": true, "message": "Invalid testId on search filter!"})

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

		return
	}

	//Create Question and Options for the test
	createQuestions, createErr := qc.questionService.CreateQuestion(reqBody.Question)
	if createErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": createErr})
		return
	}

	//Create Test in firestore using the Test created in SQL
	fsQuestionPayload := make(map[string]interface{})
	marshalledQuestion, _ := json.Marshal(createQuestions)
	json.Unmarshal(marshalledQuestion, &fsQuestionPayload)

	fsCreatedQuestion, fsQuestionDetails := qc.fireStoreService.SaveOrUpdateEntityWithId("Tests/"+testIdParam+"/Questions", createQuestions.ID, fsQuestionPayload)

	if fsQuestionDetails != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": fsQuestionDetails})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Questions and answers Created!", "data": fsCreatedQuestion})

}
