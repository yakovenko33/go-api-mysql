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
)

var (
	DB   *gorm.DB
	once sync.Once
)

func InitDBClient() {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Print(".env file not found, trying to get from environment")
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

		DB, err = gorm.Open(mysql.Open(dns), &gorm.Config{})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(err)

		sqlDB, err := DB.DB()
		if err != nil {
			log.Fatal(err) // Обработка ошибки
		}
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetConnMaxLifetime(60 * time.Second)
	})
}
