package controllers

import (
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/constants"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/dtos"
	"github/abinav-07/mcq-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	userService     services.UserService
	roleService     services.RoleService
	firebaseService services.FirebaseService
}

// Constructor
func NewAuthController(
	userServices services.UserService,
	roleService services.RoleService,
	firebaseService services.FirebaseService,
) AuthController {
	return AuthController{
		userService:     userServices,
		roleService:     roleService,
		firebaseService: firebaseService,
	}
}

// Register New User
func (uc AuthController) RegisterUser(c *gin.Context) {

	//TO:DO: Validations

	//Empty Struct
	reqBody := struct{ models.User }{}

	//Bind Body to struct
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	//Check if role exists or not
	getRole, getRoleErr := uc.roleService.GetById(reqBody.RoleID)
	if getRoleErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": "Role Not Found!"})
		return
	}

	//Encrypt Password
	hashedPassword, hashErr := utils.HashPassword(reqBody.Password)

	if hashErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": "Error while hashing password!"})
		return
	}

	//Overwrite req body password
	reqBody.Password = hashedPassword

	//Create User
	trx := c.MustGet(constants.DBTransaction).(*gorm.DB)

	claimData := dtos.UserClaimMetaData{
		UserRole: getRole.Role,
	}
	getTrx, getTrxErr := uc.userService.WithTrx(trx)
	if getTrxErr != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": getTrxErr})
		return
	}
	createdUser, createUserErr := getTrx.CreateUserWithFB(reqBody.User, claimData)

	//This also check for duplicate errors
	if createUserErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": createUserErr.Error()})
		return
	}

	validJWT, jwtErr := services.NewJWTAuthService().GenerateJWT(createdUser.ID)
	if jwtErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": jwtErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"msg": "Created User", "data": createdUser, "token": validJWT})
}

func (ac AuthController) LoginUser(c *gin.Context) {

	//Empty Struct
	reqBody := struct{ models.User }{}

	//Body to Struct
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	//Get User
	getUser, getUserErr := ac.userService.GetByEmail(reqBody.Email)

	//Check if user not found
	if getUserErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": getUserErr.Error()})
		return
	}

	if okBool := utils.CheckPasswordHash(reqBody.Password, getUser.Password); !okBool {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": "Incorrect Password"})
		return
	}

	//Get User from firebase to send FirebaseUID on response
	getFBUser, _ := ac.firebaseService.GetUserByEmail(getUser.Email)

	if getFBUser != nil {
		getUser.FirebaseUID = getFBUser.UID
	}

	//Generate Token
	validJWT, jwtErr := services.NewJWTAuthService().GenerateJWT(getUser.ID)

	if jwtErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": jwtErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": getUser, "token": validJWT})
}
