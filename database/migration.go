package database

import (
	"fmt"
	"github/abinav-07/mcq-test/infrastructure"

	"github.com/golang-migrate/migrate/v4"
	// Required for migrations to run
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"go.uber.org/fx"
)

// Migration struct
type Migration struct {
	env infrastructure.Env
}

// Return New Migration Struct
func NewMigrations(
	env infrastructure.Env,
) Migration {
	return Migration{
		env: env,
	}
}

// Migrate -> migrates all table
func (m Migration) Migrate() {
	USER := m.env.DBUsername
	PASS := m.env.DBPassword
	HOST := m.env.DBHost
	PORT := m.env.DBPort
	DBNAME := m.env.DBName

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&charset=utf8mb4&parseTime=True&loc=Local", USER, PASS, HOST, PORT, DBNAME)

	migrations, _ := migrate.New("file://./database/migration/", "mysql://"+dsn)

	//Steps >0 Up migration <0 Down migrations
	migrations.Steps(1000)
	// if err != nil {
	// 	log.Fatal("Error in migration: ", err.Error())
	// }
}

var Module = fx.Options(fx.Provide(NewMigrations))
