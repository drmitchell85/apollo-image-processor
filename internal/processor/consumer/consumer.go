package consumer

import (
	"apollo-image-processor/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type ConsumerQueue interface {
	ConsumeMessage() error
}

type consumerQueue struct {
	rmqpool *sync.Pool
}

func NewConsumerQueue(rmqpool *sync.Pool) ConsumerQueue {
	return &consumerQueue{
		rmqpool: rmqpool,
	}
}

func (c *consumerQueue) ConsumeMessage() error {

	rmqChan := c.rmqpool.Get().(*amqp.Channel)
	defer rmqChan.Close()

	queue, err := rmqChan.QueueDeclare(
		models.QueueName, // name
		false,            // durable
		false,            // auto delete
		false,            // exclusive
		false,            // no wait
		nil,              // args
	)
	if err != nil {
		return fmt.Errorf("error delcaring queue: %w", err)
	}

	msgs, err := rmqChan.Consume(
		models.QueueName, // queue
		"",               // consumer
		true,             // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              //args
	)
	if err != nil {
		return fmt.Errorf("error getting messages: %w", err)
	}

	forever := make(chan bool)
	go func() {
		for msg := range msgs {

			var batchMessage models.BatchMessage
			err = json.Unmarshal(msg.Body, &batchMessage)
			if err != nil {
				log.Printf("error getting messages: %s", err)
			}
			log.Printf("batchMessage: \n %+v", batchMessage)

		}
	}()

	log.Println("Waiting for messages...")

	log.Printf("queue status: %+v", queue)
	<-forever

	return nil
}
