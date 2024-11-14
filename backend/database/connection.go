package database

import (
	"addressBook/models/contact"
	"fmt"
	"log"
	"os"
	"os/user"
	"syscall"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectionDB() (*gorm.DB, error) {

	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		logger.Fatal("Failed to create a logger : ", zap.Error(err))
	}
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return nil, err
	}

	// Access environment variables
	// host := os.Getenv("MYSQL_HOST")
	host, found := syscall.Getenv("MYSQL_HOST")
	if !found {
		log.Println(found)
	}
	port := os.Getenv("MYSQL_PORT")
	username := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	database := os.Getenv("MYSQL_DATABASE")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True", username, password, host, port, database)
	// dsn := "user1:pass1@tcp(localhost:3306)/address_book?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Error("Failed to connect to the database :", zap.Error(err))
	}
	err = db.AutoMigrate(&user.User{}, &contact.Contact{})
	
	if err != nil {

		logger.Error("Failed to create Table :", zap.Error(err))
		return nil, err
	}
	fmt.Println("Database connection Successfully")

	return db, nil

}
