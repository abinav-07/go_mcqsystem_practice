package services

import (
	"context"
	"errors"
	"github/abinav-07/mcq-test/constants"
	"github/abinav-07/mcq-test/database/models"
	"github/abinav-07/mcq-test/dtos"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type FirebaseService struct {
	client *auth.Client
}

// NewFireBaseService constructor
func NewFirebaseService(client *auth.Client) FirebaseService {
	return FirebaseService{
		client: client,
	}
}

// Update Customer User claims
func GetUpdatedUserClaim(fb_user *auth.UserRecord, claimData dtos.UserClaimMetaData) gin.H {
	customClaim := fb_user.CustomClaims

	if len(customClaim) == 0 {
		customClaim = gin.H{}
		customClaim[constants.IsAdmin] = false
		customClaim[constants.IsUser] = true
		customClaim["role"] = constants.RoleUser
	}

	if claimData.UserRole == constants.RoleAdmin {
		customClaim[constants.IsAdmin] = true
		customClaim[constants.IsUser] = false
		customClaim["role"] = constants.RoleAdmin
	}

	customClaim["userId"] = claimData.UserID
	customClaim["firebase_uid"] = fb_user.UID
	customClaim["email"] = fb_user.Email
	customClaim["email_verified"] = fb_user.EmailVerified

	return customClaim
}

// Get or create user With given role
func (fb *FirebaseService) GetCreateOrUpdateFirebaseUser(user *models.User, claimData dtos.UserClaimMetaData) (string, error) {
	params := (&auth.UserToCreate{}).Email(user.Email).Password(user.Password)

	getUser, getUserErr := fb.GetUserByEmail(user.Email)
	customClaim := gin.H{}

	if getUser != nil && getUser.CustomClaims[claimData.UserRole] != true {
		customClaim = GetUpdatedUserClaim(getUser, claimData)
	} else if getUser == nil {
		getUser, getUserErr = fb.client.CreateUser(context.Background(), params)
		if getUserErr != nil {
			return "", getUserErr
		}

		customClaim = GetUpdatedUserClaim(getUser, claimData)
	} else {
		getUserErr = errors.New("User with the role " + claimData.UserRole + " already exists")
	}

	if getUserErr != nil {
		return "", getUserErr
	}

	//SetClaim returns only error if any
	getUserErr = fb.SetClaim(getUser.UID, customClaim)

	return getUser.UID, getUserErr
}

// Create Admin User
func (fb *FirebaseService) CreateAdminUser(user *models.User) error {
	claimData := dtos.UserClaimMetaData{
		UserID:   user.ID,
		IsAdmin:  true,
		UserRole: constants.RoleAdmin,
	}

	_, err := fb.GetCreateOrUpdateFirebaseUser(user, claimData)
	if err != nil {
		return err
	}
	return err
}

// VerifyToken ->  verify passed firebase id token
func (fb *FirebaseService) VerifyToken(idToken string) (*auth.Token, error) {
	token, err := fb.client.VerifyIDToken(context.Background(), idToken)

	return token, err
}

// Get the user data with specified email
func (fb *FirebaseService) GetUserByEmail(email string) (*auth.UserRecord, error) {
	user, err := fb.client.GetUserByEmail(context.Background(), email)
	return user, err
}

// SetClaim set's claim to firebase user
func (fb *FirebaseService) SetClaim(uid string, claims gin.H) error {

	err := fb.client.SetCustomUserClaims(context.Background(), uid, claims)
	return err
}
