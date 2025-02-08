package worker

import (
	"apollo-image-processor/internal/models"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

func ImgWorker(jobsChan <-chan amqp.Delivery, resChan chan<- amqp.Delivery) {

	for job := range jobsChan {
		time.Sleep(time.Second * 10)
		resChan <- job
	}
}

func ResWorker(resChan <-chan amqp.Delivery, rmqChan *amqp.Channel) {

	for res := range resChan {

		var batchMessage models.BatchMessage
		err := json.Unmarshal(res.Body, &batchMessage)
		if err != nil {
			log.Printf("error unmarshalling res: %s", err)
		}

		rmqChan.Ack(res.DeliveryTag, false)
	}
}
