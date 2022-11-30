package controllers

import (
	"errors"
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/constants"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/infrastructure"
	"net/http"
	"sort"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/go-cmp/cmp"
	"gorm.io/gorm"
)

type UserController struct {
	env                   infrastructure.Env
	userService           services.UserService
	roleService           services.RoleService
	testService           services.TestService
	questionService       services.QuestionService
	userTestReportService services.UserTestReportService
	firebaseService       services.FirebaseService
}

type UserTestAnswers struct {
	QuestionID uint `json:"question_id"`
	OptionID   uint `json:"option_id"`
}

// Constructor
func NewUserController(
	env infrastructure.Env,
	userService services.UserService,
	roleService services.RoleService,
	testService services.TestService,
	questionService services.QuestionService,
	userTestReportService services.UserTestReportService,
	firebaseService services.FirebaseService,
) UserController {
	return UserController{
		env:                   env,
		roleService:           roleService,
		userService:           userService,
		testService:           testService,
		questionService:       questionService,
		userTestReportService: userTestReportService,
		firebaseService:       firebaseService,
	}
}

func (uc UserController) GetAvailableTests(ctx *gin.Context) {
	query := struct{ models.Test }{}
	query.Test.IsAvailable = true
	//Get Available tests
	getTests, getTestErr := uc.testService.GetTestByQuery(query.Test)
	if getTestErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": getTestErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Test Created!", "data": getTests})
}

func (uc UserController) CreateTestReport(ctx *gin.Context) {

	//Get userid from context
	userId := uint(ctx.GetFloat64("UserId"))

	// Get body
	reqBody := struct {
		TestID      uint              `json:"test_id"`
		UserAnswers []UserTestAnswers `json:"user_answers"`
	}{}

	//Bind Body to struct
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	//TODO: User answers validations

	//Check if report already generated
	getDuplicateReport, getDuplicateErr := uc.userTestReportService.GetTestReport(userId, reqBody.TestID)
	if getDuplicateReport != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": "You have already submitted the test!"})
		return
	}
	if getDuplicateErr != nil && !errors.Is(getDuplicateErr, gorm.ErrRecordNotFound) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": getDuplicateErr.Error()})
		return
	}

	//Get Correct answers
	testCorrectAns, getAnsErr := uc.questionService.GetCorrectAnswers(reqBody.TestID)
	if getAnsErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": getAnsErr})
		return
	}

	//Check if all test questions are passed in body
	if len(testCorrectAns) != len(reqBody.UserAnswers) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": "All test questions not provided"})
		return
	}

	//Conversion of data from DB for comparison
	correctAnswersStruct := make([]UserTestAnswers, 0)
	for _, correctAnsVal := range testCorrectAns {
		optionId, _ := strconv.ParseUint(correctAnsVal.OptionId, 10, 32)
		uintOptionID := uint(optionId)
		correctAnswersStruct = append(correctAnswersStruct, UserTestAnswers{
			QuestionID: correctAnsVal.QuestionID,
			OptionID:   uintOptionID,
		})
	}

	//Sorting for correct check
	sort.Slice(correctAnswersStruct, func(i, j int) bool {
		return correctAnswersStruct[i].QuestionID < correctAnswersStruct[j].QuestionID
	})

	sort.Slice(reqBody.UserAnswers, func(i, j int) bool {
		return reqBody.UserAnswers[i].QuestionID < reqBody.UserAnswers[j].QuestionID
	})

	//Compare two structs: Body Value == Db question correct options
	has_passed := cmp.Equal(correctAnswersStruct, reqBody.UserAnswers)

	userTestReport := struct{ models.UserTestReport }{}
	userTestReport.HasPassed = has_passed
	userTestReport.TestId = reqBody.TestID
	userTestReport.UserId = userId

	createdReport, createErr := uc.userTestReportService.Create(userTestReport.UserTestReport)
	if createErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": createErr})
		return
	}

	// Custom Response payload
	payload := map[string]interface{}{
		"reportId":  createdReport.ID,
		"hasPassed": createdReport.HasPassed,
		"testId":    createdReport.TestId,
		"userId":    createdReport.UserId,
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Report Generated!", "data": payload})
}

// Deletes User From DB and FireBase Auth using Id
func (uc UserController) DeleteUserById(ctx *gin.Context) {
	//Assigned test id
	userIdParam := ctx.Param("userId")
	userId, _ := strconv.ParseUint(userIdParam, 10, 32)
	userIdParamUint := uint(userId)

	//Create Test
	trx := ctx.MustGet(constants.DBTransaction).(*gorm.DB)
	getTrx, getTrxErr := uc.userService.WithTrx(trx)

	if getTrxErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": getTrxErr})
		return
	}

	deletedUser, deleteUserErr := getTrx.DeleteById(userIdParamUint)
	if deleteUserErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": deleteUserErr})
		return
	}

	fbErr := uc.firebaseService.DeleteUser(deletedUser.Email)
	if fbErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": fbErr})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "User Deleted!", "data": "Deleted User" + deletedUser.Email})
}
