package worker

import (
	"context"
	"encoding/json"
	"log"
	"time"

	walletclient "github.com/k1ngalph0x/payflow/client/wallet"
	"github.com/k1ngalph0x/payflow/payment-service/config"
	"github.com/k1ngalph0x/payflow/payment-service/internal/events"
	"github.com/k1ngalph0x/payflow/payment-service/models"
	walletpb "github.com/k1ngalph0x/payflow/wallet-service/proto"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StartSettlementWorker(db *gorm.DB, cfg *config.Config, wc *walletclient.WalletClient, rabbitURL string) {
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("settlement worker: amqp dial: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("settlement worker: amqp channel: %v", err)
	}

	publisher, err := events.NewPublisher(rabbitURL)
	if err != nil {
		log.Fatalf("settlement worker: publisher: %v", err)
	}

	msgs, err := ch.Consume("payment.created", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("settlement worker: consume: %v", err)
	}

	for msg := range msgs {
		var event struct {
			Reference string `json:"reference"`
		}
		err = json.Unmarshal(msg.Body, &event); 
		if err != nil {
			msg.Nack(false, false)
			continue
		}

		if err := processSettlement(db, cfg, wc, publisher, event.Reference); err != nil {
			msg.Nack(false, false)
			continue
		}
		msg.Ack(false)
	}
}

func processSettlement(db *gorm.DB, cfg *config.Config, wc *walletclient.WalletClient, pub *events.Publisher, ref string) error {
	result := db.Model(&models.Payment{}).Where("reference = ? AND status = ?", ref, models.PaymentStatusCreated).Update("status", models.PaymentStatusProcessing)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return nil
	}

	var payment models.Payment
	err := db.Where("reference = ?", ref).First(&payment).Error; 
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = wc.Client.Debit(ctx, &walletpb.DebitRequest{
		UserId:    payment.UserID,
		Amount:    payment.Amount,
		Reference: ref,
	})
	if err != nil {
		db.Model(&models.Payment{}).Where("reference = ?", ref).Update("status", models.PaymentStatusFailed)
		return err
	}

	_, err = wc.Client.Credit(ctx, &walletpb.CreditRequest{
		UserId:    cfg.PLATFORM.PlatformUserID,
		Amount:    payment.Amount,
		Reference: ref,
	})
	if err != nil {
		db.Model(&models.Payment{}).Where("reference = ?", ref).Update("status", models.PaymentStatusFailed)
		return err
	}

	result = db.Model(&models.Payment{}).Clauses(clause.OnConflict{DoNothing: true}).Where("reference = ?", ref).Update("status", models.PaymentStatusFundsCaptured)
	if result.Error != nil {
		return err
	}

	return pub.Publish("payment.captured", map[string]string{"reference": ref})
}