package router

import (
	"github/abinav-07/mcq-test/infrastructure"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

type Router struct {
	Gin *gin.Engine
}

// New Router: Gin Router initializer
func NewRouter(env infrastructure.Env) Router {
	httpRouter := gin.Default()

	httpRouter.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS", "DELETE"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	httpRouter.GET("/is-working", func(ctx *gin.Context) {
		log.Print("is, working")
		ctx.JSON(http.StatusOK, gin.H{"msg": "Router is Working"})
	})
	return Router{
		Gin: httpRouter,
	}
}

var Module = fx.Options(
	fx.Provide(NewRouter),
)
