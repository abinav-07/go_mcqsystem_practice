package bootstrap

import (
	"context"
	"github/abinav-07/mcq-test/api/controllers"
	"github/abinav-07/mcq-test/api/middlewares"
	"github/abinav-07/mcq-test/api/routes"
	"github/abinav-07/mcq-test/api/services"
	"github/abinav-07/mcq-test/database"
	"github/abinav-07/mcq-test/infrastructure"
	"github/abinav-07/mcq-test/router"
	"log"

	"go.uber.org/fx"
)

var Module = fx.Options(
	router.Module,
	routes.Module,
	middlewares.Module,
	services.Module,
	controllers.Module,
	database.Module,
	infrastructure.Module,
	fx.Invoke(bootstrap),
)

// Arguments: LifeCycle Events for the application, ENV struct: values are assigned by using the fx.provide NewEnv module beforehand above, Database struct
func bootstrap(
	lifecycle fx.Lifecycle,
	env infrastructure.Env,
	appServer router.Router,
	routes routes.Routes,
	migrations database.Migration,
	database infrastructure.Database,
) {
	appStop := func(context.Context) error {
		log.Print("Stopping")
		//Setting up connection pool
		conn, _ := database.DB.DB()
		conn.Close()
		return nil
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			log.Print("Starting Go Application")
			log.Print("------------------------")

			//Starting application different thread
			go func() {
				log.Print("Migrating DB schema...")

				migrations.Migrate()
				routes.Setup()
				if env.ServerPort == "" {
					appServer.Gin.Run(":5000")
				} else {
					appServer.Gin.Run(":" + env.ServerPort)
				}
			}()
			return nil
		},
		OnStop: appStop,
	})

}
