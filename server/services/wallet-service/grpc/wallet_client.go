package grpcclient

import (
	"context"
	"time"

	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type WalletClient struct{
	Client walletpb.WalletServiceClient
}

func NewWalletClient(addr string)(*WalletClient, error){
	conn, err := grpc.NewClient(
		addr, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil{
		return nil, err
	}

	return &WalletClient{
		Client: walletpb.NewWalletServiceClient(conn),
	}, nil
}


func (w *WalletClient) CreateWallet(userID string) error{
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	_, err := w.Client.CreateWallet(ctx, &walletpb.CreateWalletRequest{
		UserId: userID,
	})

	return err
}

