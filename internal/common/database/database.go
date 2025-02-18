package database

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once         sync.Once
	errorConnect error
	db           *gorm.DB
)

func ProvideDBConnection() (*gorm.DB, error) {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			errorConnect = err
		}

		var (
			host     = os.Getenv("DB_HOST")
			port     = os.Getenv("DB_PORT")
			user     = os.Getenv("DB_USER")
			password = os.Getenv("DB_PASSWORD")
			dbname   = os.Getenv("DB_NAME")
		)

		dns := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			user,
			password,
			host,
			port,
			dbname,
		)

		db, err = gorm.Open(mysql.Open(dns), &gorm.Config{
			Logger: logger.Discard,
		})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(err)

		sqlDB, err := db.DB()
		if err != nil {
			errorConnect = err
			log.Fatalf("listen: %s\n", errorConnect)
		}
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(60 * time.Second)
	})

	return db, errorConnect
}
