package connection

import (
	"errors"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func tryToConnect() (*gorm.DB, error) {
	dsn := os.Getenv("DSN")
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}

func GetConnection() (*gorm.DB, error) {
	maxTries := 10
	counter := 0

	for {
		counter++
		db, err := tryToConnect()
		if err != nil && counter < maxTries {
			log.Println("try to connect with database...")
			time.Sleep(10 * time.Second)
			continue
		}

		if counter >= maxTries {
			log.Println("database connection timed out...")
			return nil, errors.New("timed out")
		}

		log.Println("connected with database")
		return db, nil

	}
}
