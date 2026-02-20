package api

import (
	"context"
	"errors"

	"github.com/k1ngalph0x/payflow/wallet-service/config"
	"github.com/k1ngalph0x/payflow/wallet-service/models"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletHandler struct {
	DB     *gorm.DB
	Config *config.Config
	walletpb.UnimplementedWalletServiceServer
}

func NewWalletHandler(db  *gorm.DB, config *config.Config) *WalletHandler {
	return &WalletHandler{DB: db, Config: config}
}


func (h *WalletHandler) CreateWallet(ctx  context.Context, req *walletpb.CreateWalletRequest)(*walletpb.CreateWalletResponse, error){
	wallet := models.Wallet{UserID: req.UserId}
	result :=  h.DB.WithContext(ctx).Create(&wallet)
	if result.Error != nil{
		return nil, status.Errorf(codes.Internal, "failed to create wallet: %v", result.Error)
	}
	return &walletpb.CreateWalletResponse{WalletId: wallet.ID}, nil
}

func(h *WalletHandler) GetBalance(ctx context.Context, req *walletpb.GetBalanceRequest) (*walletpb.GetBalanceResponse, error) {
	var wallet models.Wallet
	result := h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).First(&wallet)
	if result.Error != nil{
		return nil, status.Errorf(codes.Internal, "failed to fetch balance: %v", result.Error)
	}

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, status.Errorf(codes.NotFound, "wallet not found")
	}

	return &walletpb.GetBalanceResponse{Balance: wallet.Balance}, nil
}

func (h *WalletHandler) GetTransactions(ctx context.Context, req *walletpb.GetTransactionsRequest) (*walletpb.GetTransactionsResponse, error){

	var transactions []models.Transaction
	result :=  h.DB.WithContext(ctx).Where("user_id = ?", req.UserId).Order("created_at DESC").Limit(int(req.Limit)).Offset(int(req.Offset)).Find(&transactions)
	if result.Error != nil{
		return nil, status.Errorf(codes.Internal, "failed to fetch transactions: %v", result.Error)
	}

	walletTransactions := make([]*walletpb.Transaction, len(transactions))
	for i, t:= range transactions{
		walletTransactions[i] = &walletpb.Transaction{
			Id:        t.ID,
			Type:      string(t.Type),
			Amount:    t.Amount,
			Reference: t.Reference,
			Status:    string(t.Status),
			CreatedAt: t.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}
	return &walletpb.GetTransactionsResponse{Transactions: walletTransactions}, nil
}

func (h *WalletHandler) Debit(ctx context.Context, req *walletpb.DebitRequest) (*walletpb.TransactionResponse, error) {
	var transactionID string
	var responseStatus string

	err := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", req.UserId).First(&wallet).Error
		if err != nil {
			return status.Errorf(codes.NotFound, "wallet not found")
		}

		var existing models.Transaction
		err = tx.Where("wallet_id = ? AND reference = ? AND type = ?",wallet.ID, req.Reference, models.TransactionTypeDebit).First(&existing).Error
		if err == nil {
			transactionID = existing.ID
			responseStatus = string(models.TransactionStatusSuccess)
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.Internal, "something went wrong: %v", err)
		}
		if wallet.Balance < req.Amount {
			responseStatus = string(models.TransactionStatusInsufficientFunds)
			return nil
		}

		txn := models.Transaction{
			WalletID:  wallet.ID,
			UserID:    req.UserId,
			Type:      models.TransactionTypeDebit,
			Amount:    req.Amount,
			Reference: req.Reference,
			Status:    models.TransactionStatusSuccess,
		}
		if err := tx.Create(&txn).Error; err != nil {
			return status.Errorf(codes.Internal, "failed to insert transaction: %v", err)
		}

		err = tx.Model(&wallet).Update("balance", gorm.Expr("balance - ?", req.Amount)).Error; 
		if err != nil {
			return status.Errorf(codes.Internal, "failed to update balance: %v", err)
		}

		transactionID = txn.ID
		responseStatus = string(models.TransactionStatusSuccess)
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &walletpb.TransactionResponse{TransactionId: transactionID, Status: responseStatus}, nil
}


func (h *WalletHandler) Credit(ctx context.Context, req *walletpb.CreditRequest) (*walletpb.TransactionResponse, error) {
	var transactionID string

	err := h.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var wallet models.Wallet
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("user_id = ?", req.UserId).First(&wallet).Error; 
		if err != nil {
			return status.Errorf(codes.NotFound, "wallet not found")
		}

		var existing models.Transaction
		err = tx.Where("wallet_id = ? AND reference = ? AND type = ?",wallet.ID, req.Reference, models.TransactionTypeCredit).First(&existing).Error
		if err == nil {
			transactionID = existing.ID
			return nil
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return status.Errorf(codes.Internal, "something went wrong: %v", err)
		}

		txn := models.Transaction{
			WalletID:  wallet.ID,
			UserID:    req.UserId,
			Type:      models.TransactionTypeCredit,
			Amount:    req.Amount,
			Reference: req.Reference,
			Status:    models.TransactionStatusSuccess,
		}
		if err := tx.Create(&txn).Error; err != nil {
			return status.Errorf(codes.Internal, "failed to insert transaction: %v", err)
		}

		if err := tx.Model(&wallet).Update("balance", gorm.Expr("balance + ?", req.Amount)).Error; err != nil {
			return status.Errorf(codes.Internal, "failed to update balance: %v", err)
		}

		transactionID = txn.ID
		return nil
	})

	if err != nil {
		return nil, err
	}
	return &walletpb.TransactionResponse{TransactionId: transactionID, Status: string(models.TransactionStatusSuccess)}, nil
}