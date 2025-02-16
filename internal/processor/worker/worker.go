package worker

import (
	"apollo-image-processor/internal/models"
	"apollo-image-processor/internal/processor/imager"
	procrepository "apollo-image-processor/internal/processor/repository"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/streadway/amqp"
)

// TODO address if the continue will cause issues with the message buffer...
func ImgWorker(jobsChan <-chan amqp.Delivery, resChan chan<- amqp.Delivery, errChan chan<- error, db *sql.DB) {

	for job := range jobsChan {

		var batchMessage models.BatchMessage
		err := json.Unmarshal(job.Body, &batchMessage)
		if err != nil {
			log.Printf("error unmarshalling res: %s", err)
		}

		srcimage, err := procrepository.GetImage(batchMessage.Imageid, db)
		if err != nil {
			errChan <- err
			continue
		}

		procImage, err := imager.ProcessImageBW(srcimage)
		if err != nil {
			errChan <- err
			continue
		}

		// fmt.Printf("\nimage %s processed", batchMessage.Imageid)

		err = procrepository.InsertImage(batchMessage.Imageid, procImage, db)
		if err != nil {
			errChan <- err
			continue
		}

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

func ErrWorker(errChan <-chan error) {

	for err := range errChan {

		log.Printf("error: %s", err)

	}
}
