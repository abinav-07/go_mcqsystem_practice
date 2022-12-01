package controllers

import (
	"fmt"
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/constants"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/dtos"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController struct {
	userService     services.UserService
	roleService     services.RoleService
	firebaseService services.FirebaseService
}

// Constructor
func NewAdminController(
	userService services.UserService,
	roleService services.RoleService,
	firebaseService services.FirebaseService,
) AdminController {
	return AdminController{
		userService:     userService,
		roleService:     roleService,
		firebaseService: firebaseService,
	}
}

// Update User Data and Claims
func (ac AdminController) UpdateUser(ctx *gin.Context) {
	//Empty Struct
	reqBody := struct{ models.User }{}

	//Assigned user id
	userIdParam := ctx.Param("userId")
	userId, _ := strconv.ParseUint(userIdParam, 10, 32)
	userIdParamUint := uint(userId)

	//Bind Body to struct
	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": err.Error()})
		return
	}

	fmt.Println("Role Id", reqBody.RoleID)

	claimData := dtos.UserClaimMetaData{}

	if reqBody.RoleID != 0 {
		// Check if role exists or not
		getRole, getRoleErr := ac.roleService.GetById(reqBody.RoleID)
		if getRoleErr != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": "Role Not Found!"})
			return
		}

		claimData.UserRole = getRole.Role
	}

	//Create User
	trx := ctx.MustGet(constants.DBTransaction).(*gorm.DB)
	getTrx, getTrxErr := ac.userService.WithTrx(trx)
	if getTrxErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": getTrxErr})
		return
	}

	updatedUser, updateUserErr := getTrx.UpdateOneUserWithFB(userIdParamUint, reqBody.User, claimData)

	//This also check for duplicate errors
	if updateUserErr != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": true, "message": updateUserErr.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "Updated User", "data": updatedUser})
}
