package api

import (
	"context"
	"database/sql"
	"time"

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

func(h *WalletHandler) GetTransactions(ctx context.Context, req *walletpb.GetTransactionsRequest)(*walletpb.GetTransactionsResponse, error){

	var transactions[]*walletpb.Transaction

	query := `SELECT id, type, amount, reference, status, created_at
		FROM payflow_wallet_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := h.DB.Query(query, req.UserId, req.Limit, req.Offset)

	if err != nil{
		return nil, err
	}

	defer rows.Close()

	for rows.Next(){
		var txn walletpb.Transaction
		var createdAt time.Time

		err := rows.Scan(
			&txn.Id,
			&txn.Type,
			&txn.Amount,
			&txn.Reference,
			&txn.Status,	
			&createdAt,
		)

		if err != nil{
			return nil, err
		}

		txn.CreatedAt = createdAt.Format(time.RFC3339)
		transactions = append(transactions, &txn)
	}

	return &walletpb.GetTransactionsResponse{Transactions: transactions}, nil
}

func (h *WalletHandler) Debit(ctx context.Context, req *walletpb.DebitRequest)(*walletpb.TransactionResponse, error){
	var walletID string
	var balance float64
	var txnId string
	tx, err := h.DB.BeginTx(ctx, nil)
	if err != nil{
		return nil, err
	}
	defer tx.Rollback()

	query := `SELECT id, balance FROM payflow_wallets WHERE user_id = $1 FOR UPDATE`

	err = tx.QueryRow(query, req.UserId).Scan(&walletID, &balance)

	if err != nil{
		return nil, err
	}

	if balance < req.Amount{
		return &walletpb.TransactionResponse{
			Status: "INSUFFICIENT_FUNDS",

		}, nil
	}

	txnQuery := `INSERT INTO payflow_wallet_transactions
	(wallet_id, user_id, type, amount, reference, status)
	VALUES ($1, $2, 'DEBIT', $3, $4, 'SUCCESS') RETURNING id
	`
	err = tx.QueryRow(txnQuery, walletID, req.UserId, req.Amount, req.Reference).Scan(&txnId)
	if err!= nil{
		return nil, err
	}

	updateTxnQuery := `UPDATE payflow_wallets 
	SET balance = balance - $1, updated_at = NOW()
	WHERE id = $2`
	_, err = tx.Exec(updateTxnQuery, req.Amount, walletID)

	if err!= nil{
		return nil, err
	}

	
	err = tx.Commit()

	if err != nil{
		return nil, err
	}

	return &walletpb.TransactionResponse{
		TransactionId: txnId,
		Status: "SUCCESS",
	}, nil
}

func(h *WalletHandler) Credit(ctx context.Context, req *walletpb.CreditRequest)(*walletpb.TransactionResponse, error){
	var walletID string
	var txnId string
	tx, err := h.DB.BeginTx(ctx, nil)
	if err!= nil{
		return nil, err
	}

	defer tx.Rollback()

	query := `SELECT id FROM payflow_wallets WHERE user_id = $1 FOR UPDATE`
	
	err = tx.QueryRow(query, req.UserId).Scan(&walletID)

	if err!= nil{
		return nil, err
	}

	insertQuery := `
	INSERT INTO payflow_wallet_transactions
	(wallet_id, user_id, type, amount, reference, status)
	VALUES ($1, $2, 'CREDIT', $3, $4, 'SUCCESS') RETURNING id
	`

	err = tx.QueryRow(insertQuery, walletID, req.UserId, req.Amount, req.Reference).Scan(&txnId)
	
	if err!= nil{
		return nil, err
	}

	updateTxnQuery := `UPDATE payflow_wallets 
	SET balance = balance + $1, updated_at = NOW()
	WHERE id = $2`
	_, err = tx.Exec(updateTxnQuery, req.Amount, walletID)

	if err!= nil{
		return nil, err
	}

	err = tx.Commit()

	if err != nil{
		return nil, err
	}

	return &walletpb.TransactionResponse{
		TransactionId: txnId,
		Status: "SUCCESS",
	}, nil
}