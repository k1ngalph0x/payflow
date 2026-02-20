package db

import (
	"fmt"

	"github.com/k1ngalph0x/payflow/identity-service/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ConnectDB() (*gorm.DB, error) {
	config, err := config.LoadConfig()

	if err != nil{
		return nil, err
	}

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Password, config.DB.Dbname)

	db, err := gorm.Open(postgres.Open(conn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil{
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	fmt.Println("Successfully connected to database")

	return db, nil

}