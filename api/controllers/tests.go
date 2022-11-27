package controllers

import (
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TestController struct {
	testService services.TestService
}

// Construct
func NewTestController(
	testService services.TestService,
) TestController {
	return TestController{
		testService: testService,
	}
}

// Create Tests
func (tc TestController) CreateTests(ctx *gin.Context) {
	//TO:DO: Validations

	//Empty Struct
	reqBody := struct{ models.Test }{}

	//Bind Body to test struct
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	reqBody.Test.IsAvailable = true

	//Check if Test already exists or not
	if _, err := tc.testService.GetByName(reqBody.Title); err == nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": "Duplicate test name found!"})
		return
	}

	createdTest, createdTestErr := tc.testService.Create(reqBody.Test)
	if createdTestErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": createdTestErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Created!", "data": createdTest})

}

func (tc TestController) GetTestDetails(ctx *gin.Context) {

	//Empty Struct
	testIdParam := ctx.Param("testId")

	if testIdParam == "" {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": true, "message": "Invalid testId on search filter!"})
		ctx.Abort()
		return
	}

	testId, _ := strconv.ParseUint(testIdParam, 10, 32)
	testIdParamUint := uint(testId)

	testDetails, testErr := tc.testService.GetById(testIdParamUint)
	if testErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, " message": testErr.Error()})
		ctx.Abort()
		return

	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Details!", "data": testDetails})

}