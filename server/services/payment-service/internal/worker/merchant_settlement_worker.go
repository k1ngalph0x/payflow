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

func StartMerchantSettlementWorker(
	db *sql.DB,
	cfg *config.Config,
	walletClient *walletclient.WalletClient,
	rabbitURL string,
) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatal(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	msgs, _ := ch.Consume(
		"payment.captured",
		"",
		false,
		false,
		false,
		false,
		nil,
	)

	for msg := range msgs {
		var event struct {
			Reference string `json:"reference"`
		}

		json.Unmarshal(msg.Body, &event)

		err := settleToMerchant(db, cfg, walletClient, event.Reference); 
		if err != nil {
			log.Println("merchant settlement failed:", err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}
}


func settleToMerchant(
	db *sql.DB,
	cfg *config.Config,  
	walletClient *walletclient.WalletClient,
	ref string,
) error {

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var payment struct {
		MerchantId string
		MerchantUserId string
		Amount     float64
		Status     string
	}

	err = tx.QueryRow(`
		SELECT merchant_id, merchant_user_id, amount, status
		FROM payflow_payments
		WHERE reference = $1
		FOR UPDATE
	`, ref).Scan(
		&payment.MerchantId,
		&payment.MerchantUserId,
		&payment.Amount,
		&payment.Status,
	)

	if err != nil {
		return err
	}

	if payment.Status != "FUNDS_CAPTURED" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = walletClient.Client.Debit(ctx, &walletpb.DebitRequest{
		UserId:   cfg.PLATFORM.PlatformUserID,
		Amount:   payment.Amount,
		Reference: ref + "_settle",
	})
	if err != nil {
		return err
	}

	_, err = walletClient.Client.Credit(ctx, &walletpb.CreditRequest{
		//UserId:    payment.MerchantId,
		UserId: payment.MerchantUserId,
		Amount:   payment.Amount,
		Reference: ref + "_settle",
	})
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		UPDATE payflow_payments
		SET status = 'SETTLED'
		WHERE reference = $1
	`, ref)
	if err != nil {
		return err
	}

	return tx.Commit()
}
