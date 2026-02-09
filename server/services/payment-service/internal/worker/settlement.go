package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	amqp "github.com/rabbitmq/amqp091-go"
)

func StartSettlementWorker(
	db *sql.DB,
	cfg *config.Config,
	walletClient *walletclient.WalletClient,
	rabbitURL string,
) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil{
		log.Fatal("payment-service/rabbitmq: - conn",err)
	}
	ch, err := conn.Channel()
	if err != nil{
		log.Fatal("payment-service/rabbitmq: - ch",err)
	}
	// var event struct{
	// 	Reference string `json:"reference"`
	// }


	msgs, _ := ch.Consume(
		"payment.created",
		"",
		false, 
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs{
		var event struct{
			Reference string `json:"reference"`
		}
		json.Unmarshal(msg.Body, &event)
		err := processSettlement(db, cfg, walletClient, event.Reference)
		if err != nil{
			log.Println("Settlement falied", err)
			msg.Nack(false, false)
			continue
		}

		msg.Ack(false)
	}
}

func processSettlement(
	db *sql.DB,
	cfg *config.Config,
	walletClient *walletclient.WalletClient,
	ref string,
)error {
	
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback() 

	var payment struct{
		UserId string
		MerchantId string
		Amount float64
		Status string
	}

	query := `
		SELECT user_id, merchant_id, amount, status
		FROM payflow_payments
		WHERE reference = $1
		FOR UPDATE
	`
	err = tx.QueryRow(query, ref).Scan(&payment.UserId, &payment.MerchantId, &payment.Amount, &payment.Status,)

	if err != nil{
		return err
	}

	if payment.Status != "CREATED"{
		return nil
	}

	updateQuery := `UPDATE payflow_payments SET status = 'PROCESSING' WHERE reference = $1`

	_, err = tx.Exec(updateQuery, ref)

	if err!= nil{
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = walletClient.Client.Debit(ctx, &walletpb.DebitRequest{
		UserId: payment.UserId,
		Amount: payment.Amount,
		Reference: ref,
	})

	if err != nil{
		db.Exec(`UPDATE payflow_payments SET status='FAILED' WHERE reference=$1`, ref)
		return err 
	}

	_, err = walletClient.Client.Credit(ctx, &walletpb.CreditRequest{
		UserId:    cfg.PLATFORM.PlatformUserID,
		Amount: payment.Amount,
		Reference: ref,
	})

	if err != nil{
		db.Exec(`UPDATE payflow_payments SET status='FAILED' WHERE reference=$1`, ref)
		return err 
	}

	_, err = db.Exec(
		`UPDATE payflow_payments SET status='FUNDS_CAPTURED' WHERE reference=$1`,
		ref,
	)

	return err

}