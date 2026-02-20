package events

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	Conn *amqp.Connection
	Channel *amqp.Channel
}



func NewPublisher(url string)(*Publisher, error){
	conn, err := amqp.Dial(url)
	queues := []string{
		"payment.created",
		"payment.captured",
	}

	if err != nil{
		log.Fatal("payment-service/rabbitmq: - conn",err)
	}

	ch, err := conn.Channel()
		if err != nil{
		log.Fatal("payment-service/rabbitmq: - ch",err)
	}
	
	for _, q := range queues {
	_, err = ch.QueueDeclare(
		q,
		true,
		false,
		false, 
		false,
		nil,
		)
		if err!=nil{
			return nil,  err
		}
	}


	if err!=nil{
		return nil,  err
	}

	return &Publisher{
		Conn:  conn,
		Channel: ch,
	}, nil
}

func (p *Publisher) Publish(queue string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil{
		return err
	}

	return p.Channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body: body,
		},
	)
}