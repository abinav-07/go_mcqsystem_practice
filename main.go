package main

import (
	"github/abinav-07/mcq-test/bootstrap"
	"log"

	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Panic("Error loading .env file: ", envErr)
	}

	fx.New(bootstrap.Module).Run()
}
