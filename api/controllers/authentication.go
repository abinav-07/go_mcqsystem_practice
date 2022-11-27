package controllers

import (
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService services.UserService
	roleService services.RoleService
}

// Constructor
func NewAuthController(
	userServices services.UserService,
	roleService services.RoleService,
) AuthController {
	return AuthController{
		userService: userServices,
		roleService: roleService,
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
	if _, err := uc.roleService.GetById(reqBody.RoleID); err != nil {
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
	createdUser, createUserErr := uc.userService.Create(reqBody.User)

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

	//Generate Token
	validJWT, jwtErr := services.NewJWTAuthService().GenerateJWT(getUser.ID)

	if jwtErr != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": true, "message": jwtErr})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": getUser, "token": validJWT})
}
