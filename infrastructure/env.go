package infrastructure

import "os"

type Env struct {
	ServerPort       string
	Environment      string
	DBUsername       string
	DBPassword       string
	DBHost           string
	DBPort           string
	DBName           string
	SecretKey        string
	FireStoreProject string
}

// Creates a new Environment
func NewEnv() Env {
	env := Env{}
	env.LoadEnv()
	return env
}

func (env *Env) LoadEnv() {
	env.ServerPort = os.Getenv("ServerPort")
	env.Environment = os.Getenv("Environment")

	env.DBUsername = os.Getenv("DBUsername")
	env.DBPassword = os.Getenv("DBPassword")
	env.DBHost = os.Getenv("DBHost")
	env.DBPort = os.Getenv("DBPort")
	env.DBName = os.Getenv("DBName")

	env.SecretKey = os.Getenv("JWTSecretKey")

	env.FireStoreProject = os.Getenv("FireStoreProject")
}
