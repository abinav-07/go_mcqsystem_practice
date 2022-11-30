package controllers

import (
	"encoding/json"
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/constants"
	"github/abinav-07/mcq-test/database/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TestController struct {
	testService      services.TestService
	firestoreService services.FireStoreService
}

// Construct
func NewTestController(
	testService services.TestService,
	firestoreService services.FireStoreService,
) TestController {
	return TestController{
		testService:      testService,
		firestoreService: firestoreService,
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

	//Create Test
	trx := ctx.MustGet(constants.DBTransaction).(*gorm.DB)
	getTrx, getTrxErr := tc.testService.WithTrx(trx)
	if getTrxErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": getTrxErr})
		return
	}
	createdTest, createdTestErr := getTrx.Create(reqBody.Test)
	if createdTestErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": createdTestErr})
		return
	}

	//Create Test in firestore using the Test created in SQL
	fsTestPayload := make(map[string]interface{})
	marshalledTest, _ := json.Marshal(createdTest)
	json.Unmarshal(marshalledTest, &fsTestPayload)
	_, fsTestErr := tc.firestoreService.SaveOrUpdateEntityWithId("Tests", createdTest.ID, fsTestPayload)

	if fsTestErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": fsTestErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Created!", "data": createdTest})

}

func (tc TestController) GetTestDetails(ctx *gin.Context) {

	//Assigned test id
	testIdParam := ctx.Param("testId")

	if testIdParam == "" {
		ctx.JSON(http.StatusForbidden,
			gin.H{"error": true, "message": "Invalid testId on search filter!"})

		return
	}

	testId, _ := strconv.ParseUint(testIdParam, 10, 32)
	testIdParamUint := uint(testId)

	//Get Test By Id From Database
	_, testErr := tc.testService.GetById(testIdParamUint)
	if testErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, " message": testErr.Error()})

		return

	}

	//Get Test By Id From Firestore
	fsTestDetails, fsTestErr := tc.firestoreService.GetEntityWithId("Tests", testIdParamUint)

	if fsTestErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, " message": fsTestErr.Error()})

		return

	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Details!", "data": fsTestDetails.Data()})

}

func (tc TestController) UpdatePartial(ctx *gin.Context) {

	//Empty Struct
	reqBody := struct{ models.Test }{}

	// Bind Body to test struct
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	//Assigned test id
	testIdParam := ctx.Param("testId")
	testId, _ := strconv.ParseUint(testIdParam, 10, 32)
	testIdParamUint := uint(testId)

	//Update Test to unavailable
	if _, err := tc.testService.UpdateOneTest(testIdParamUint, map[string]interface{}{
		"is_available": reqBody.Test.IsAvailable,
	}); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, " message": err.Error()})

		return
	}

	//For Firestore
	fsUpdateTestPayload := make(map[string]interface{})
	marshalledTest, _ := json.Marshal(reqBody.Test)
	json.Unmarshal(marshalledTest, &fsUpdateTestPayload)

	_, fsUpdatedTestErr := tc.firestoreService.SaveOrUpdateEntityWithId("Tests", testIdParamUint, fsUpdateTestPayload)
	if fsUpdatedTestErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": fsUpdatedTestErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Updated!"})
}

func (fsc TestController) DeleteCollectionDocument(ctx *gin.Context) {
	//Assigned test id
	testIdParam := ctx.Param("testId")
	testId, _ := strconv.ParseUint(testIdParam, 10, 32)
	testIdParamUint := uint(testId)

	//Create Test
	trx := ctx.MustGet(constants.DBTransaction).(*gorm.DB)
	getTrx, getTrxErr := fsc.testService.WithTrx(trx)
	if getTrxErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": getTrxErr})
		return
	}
	_, deleteTestErr := getTrx.DeleteById(testIdParamUint)
	if deleteTestErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": deleteTestErr})
		return
	}

	//Delete document from Firestore
	err := fsc.firestoreService.DeleteCollectionWithId("Tests", testIdParamUint)

	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": err})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Deleted!"})
}
