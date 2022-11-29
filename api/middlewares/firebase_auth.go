package middlewares

import (
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/constants"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/gin-gonic/gin"
)

type FirebaseAuthMW struct {
	service services.FirebaseService
}

// Construct
func NewFirebaseAuthMiddleware(
	service services.FirebaseService,
) FirebaseAuthMW {
	return FirebaseAuthMW{
		service: service,
	}
}

// Get token from custom header
func (m FirebaseAuthMW) GetTokenFromCustomHeader(ctx *gin.Context) (*auth.Token, error) {
	header := ctx.GetHeader(constants.CustomAuthorization)
	idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

	return m.service.VerifyToken(idToken)
}

// Compare roles from Claims
func (m FirebaseAuthMW) handleUserTypeVerification(ctx *gin.Context, token *auth.Token, userTypes ...string) {
	hasPermission := false

	for _, userType := range userTypes {
		authorized, ok := token.Claims[userType]

		if ok && authorized.(bool) {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Un-authorized User"})
		ctx.Abort()
		return
	}

	//Set Claims
	ctx.Set(constants.Claims, token.Claims)
}

// Handle for only Admins
func (m FirebaseAuthMW) HandleAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := m.GetTokenFromCustomHeader(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": "Un-authorized User"})
			ctx.Abort()
			return
		}

		m.handleUserTypeVerification(ctx, token, constants.IsAdmin)
		ctx.Next()

	}
}

// Handle for only users
func (m FirebaseAuthMW) HandleUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, err := m.GetTokenFromCustomHeader(ctx)

		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": true, "message": err})
			ctx.Abort()
			return
		}

		m.handleUserTypeVerification(ctx, token, constants.IsUser)
		ctx.Next()

	}
}
