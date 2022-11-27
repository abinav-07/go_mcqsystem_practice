package middlewares

import (
	"github/abinav-07/mcq-test/api/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"go.uber.org/fx"
)

type AuthMW struct {
	userService services.UserService
}

func NewAuthMiddlware(userService services.UserService) AuthMW {
	return AuthMW{
		userService: userService,
	}
}

func (a AuthMW) CheckJWT() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		const BEARER_SCHEMA = "Bearer "

		authHeader := ctx.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_SCHEMA):]

		validJWT, err := services.NewJWTAuthService().ValidateToken(tokenString)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "error": err})
		}

		if validJWT.Valid {
			claims := validJWT.Claims.(jwt.MapClaims)
			ctx.Set("UserId", claims["user_id"])

			//Check User exists and set role to context
			userId := uint(claims["user_id"].(float64))
			getUser, getUsererr := a.userService.GetById(userId)

			if getUsererr != nil {

				ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": getUsererr})
				ctx.Abort()
				return
			}

			ctx.Set("Role", getUser.Role.Role)

			ctx.Next()
		} else {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": http.StatusUnauthorized, "error": "Please enter a valid token string."})
		}

	}

}

func (a AuthMW) CheckAdmin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userRole := ctx.GetString("Role")

		if userRole != "admin" {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": true, "message": "Only admins are allowed!"})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

var Module = fx.Options(fx.Provide(NewAuthMiddlware))
