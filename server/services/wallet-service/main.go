package main

import (
	"fmt"
	"log"
	"net"

	"github.com/k1ngalph0x/payflow/wallet-service/api"
	"github.com/k1ngalph0x/payflow/wallet-service/config"
	"github.com/k1ngalph0x/payflow/wallet-service/db"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	"google.golang.org/grpc"
)

func main() {

	cfg, err := config.LoadConfig()


	if err != nil {
		log.Fatalf("Error loading config: ", err)
	}

	//Connect to db
	conn, err := db.ConnectDB()

	if err != nil {
		log.Fatalf("Error connecting to database: ", err)
	}

	defer conn.Close()	

	//Connect to port tcp
	listener, err := net.Listen("tcp", ":50051")
	if err != nil{
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	walletHandler := api.NewWalletHandler(conn, cfg)
	
	walletpb.RegisterWalletServiceServer(grpcServer, walletHandler)

	fmt.Println("Running wallet-service")

	grpcServer.Serve(listener)
}