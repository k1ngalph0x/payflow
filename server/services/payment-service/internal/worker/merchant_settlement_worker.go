package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	"github.com/k1ngalph0x/payflow/payment-service/models"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
)

func StartMerchantSettlementWorker(db *gorm.DB, cfg *config.Config, wc *walletclient.WalletClient, rabbitURL string) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("merchant settlement worker: amqp dial: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("merchant settlement worker: amqp channel: %v", err)
	}

	msgs, err := ch.Consume("payment.captured", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("merchant settlement worker: consume: %v", err)
	}

	for msg := range msgs {
		var event struct {
			Reference string `json:"reference"`
		}
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			msg.Nack(false, false)
			continue
		}

		if err := settleToMerchant(db, cfg, wc, event.Reference); err != nil {
			msg.Nack(false, true) 
			continue
		}
		msg.Ack(false)
	}
}

func settleToMerchant(db *gorm.DB, cfg *config.Config, wc *walletclient.WalletClient, ref string) error {
	var payment models.Payment
	if err := db.Where("reference = ? AND status = ?", ref, models.PaymentStatusFundsCaptured).
		First(&payment).Error; err != nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	settleRef := ref + "_settle"
	_, err := wc.Client.Debit(ctx, &walletpb.DebitRequest{
		UserId:    cfg.PLATFORM.PlatformUserID,
		Amount:    payment.Amount,
		Reference: settleRef,
	})
	if err != nil {
		return err
	}
	
	_, err = wc.Client.Credit(ctx, &walletpb.CreditRequest{
		UserId:    payment.MerchantUserID,
		Amount:    payment.Amount,
		Reference: settleRef,
	})
	if err != nil {
		return err
	}

	result := db.Model(&models.Payment{}).Where("reference = ? AND status = ?", ref, models.PaymentStatusFundsCaptured).Update("status", models.PaymentStatusSettled)
	if result.Error != nil{
		return result.Error
	}

	return nil
}