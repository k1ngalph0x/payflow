package api

import (
	"context"
	"database/sql"

	"github.com/k1ngalph0x/payflow/wallet-service/config"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
)

type WalletHandler struct {
	DB     *sql.DB
	walletpb.UnimplementedWalletServiceServer
	Config *config.Config
}


func NewWalletHandler(db *sql.DB, cfg *config.Config) *WalletHandler {
	return &WalletHandler{DB: db, Config: cfg}
}


func (h *WalletHandler) CreateWallet(ctx  context.Context, req *walletpb.CreateWalletRequest)(*walletpb.CreateWalletResponse, error){
	var walletId string
	query := `INSERT INTO payflow_wallets (user_id) VALUES ($1) RETURNING id`

	err := h.DB.QueryRow(query, req.UserId).Scan(&walletId)
	if err!=nil{
		return nil, err
	}

	return &walletpb.CreateWalletResponse{WalletId: walletId}, nil

}



func(h *WalletHandler) GetBalance(ctx context.Context, req *walletpb.GetBalanceRequest)(*walletpb.GetBalanceResponse, error){
	var balance float64
	query := `SELECT balance FROM payflow_wallets WHERE user_id = $1`

	err := h.DB.QueryRow(query, req.UserId).Scan(&balance)
	if err != nil{
		return nil, err
	}

	return &walletpb.GetBalanceResponse{Balance: balance}, nil
}