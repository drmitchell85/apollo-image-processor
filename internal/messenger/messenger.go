package messenger

import (
	"apollo-image-processor/internal/models"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type MessengerQueue interface {
	PublishMessage(string, []string) error
}

type messengerQueue struct {
	rmqpool *sync.Pool
}

func NewMessengerQueue(rmqpool *sync.Pool) MessengerQueue {
	return &messengerQueue{
		rmqpool: rmqpool,
	}
}

func (m *messengerQueue) PublishMessage(batchID string, imageIDs []string) error {

	rmqChan := m.rmqpool.Get().(*amqp.Channel)
	defer rmqChan.Close()

	queue, err := rmqChan.QueueDeclare(
		"processbatch", // name
		false,          // durable
		false,          // auto delete
		false,          // exclusive
		false,          // no wait
		nil,            // args
	)
	if err != nil {
		return fmt.Errorf("error delcaring queue: %w", err)
	}

	for i := 0; i < len(imageIDs); i++ {

		batchMessage := models.BatchMessage{
			Batchid: batchID,
			Imageid: imageIDs[i],
		}

		b, err := json.Marshal(batchMessage)
		if err != nil {
			return fmt.Errorf("error error marshalling batch message: %w", err)
		}

		err = rmqChan.Publish(
			"",         // exchange
			queue.Name, // key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				// Body:        []byte(b),
				Body: b,
			},
		)
		if err != nil {
			return fmt.Errorf("error publishing message: %w", err)
		}

	}

	log.Printf("queue status: %+v", queue)

	return nil
}
