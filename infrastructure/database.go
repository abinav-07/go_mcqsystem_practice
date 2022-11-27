package infrastructure

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database modal
type Database struct {
	DB *gorm.DB
}

// Creates Database Instance
func NewMCQTestDatabase(env Env) Database {

	//Gorm Logger to set new logs values for gorm config. Displays time when gorm was initialized
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level info: Color green (from docs)
			Colorful:      true,        // Enable color
		},
	)

	fmt.Println("Env: ", env)

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&charset=utf8mb4&parseTime=True&loc=Local", env.DBUsername, env.DBPassword, env.DBHost, env.DBPort, env.DBName)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{Logger: newLogger})
	log.Print("DB: ", db, url)

	db.Exec("CREATE DATABASE IF NOT EXISTS " + env.DBName + ";")

	if err != nil {
		log.Fatal(err)
	}

	return Database{
		DB: db,
	}
}
