package db

import (
	"database/sql"
	"fmt"

	"github.com/k1ngalph0x/payflow/auth-service/config"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error){

	config, err := config.LoadConfig()

	if err!=nil{
		return nil, err
	}

	conn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", config.DB.Host, config.DB.Port, config.DB.Username, config.DB.Password, config.DB.Dbname)

	db, err := sql.Open("postgres", conn)

	if err != nil{
		//panic(err)
		return nil, err
	}

	//defer db.Close()

	err = db.Ping()

	if err != nil{
		//panic(err)
		return nil, err
	}

	fmt.Println("Successfully connected to Database!")

	return db, nil

}