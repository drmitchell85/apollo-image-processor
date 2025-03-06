package worker

import (
	"apollo-image-processor/internal/models"
	"apollo-image-processor/internal/processor/imager"
	procrepository "apollo-image-processor/internal/processor/repository"
	"encoding/json"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

func ImgWorker(
	jobsChan <-chan amqp.Delivery,
	resChan chan<- amqp.Delivery,
	errChan chan<- error,
	shutdown <-chan struct{},
	wg *sync.WaitGroup,
	pr procrepository.ProcessorRepository,
) {

	for {
		select {
		case job := <-jobsChan:
			var batchMessage models.BatchMessage
			err := json.Unmarshal(job.Body, &batchMessage)
			if err != nil {
				log.Printf("error unmarshalling res: %s", err)
				errChan <- err
				continue
			}

			srcimage, err := pr.GetImage(batchMessage.Imageid, batchMessage.Batchid)
			if err != nil {
				errChan <- err
				continue
			}

			procImage, err := imager.ProcessImageBW(srcimage)
			if err != nil {
				errChan <- err
				continue
			}

			err = pr.InsertImage(batchMessage.Imageid, batchMessage.Batchid, procImage)
			if err != nil {
				errChan <- err
				continue
			}

			// time.Sleep(time.Second * 10)
			resChan <- job

		case <-shutdown:
			wg.Done()
			return

		default:
			// continue...
		}

	}
}

func ResWorker(resChan <-chan amqp.Delivery, shutdown <-chan struct{}, rmqChan *amqp.Channel, wg *sync.WaitGroup) {

	for {
		select {
		case res := <-resChan:
			var batchMessage models.BatchMessage
			err := json.Unmarshal(res.Body, &batchMessage)
			if err != nil {
				log.Printf("error unmarshalling res: %s", err)
			}

			rmqChan.Ack(res.DeliveryTag, false)

		case <-shutdown:
			wg.Done()
			return

		default:
			// continue...
		}
	}
}

// TODO fix to handle error messages to queue
// must updated database status
// must ack errors
func ErrWorker(errChan <-chan error, shutdown <-chan struct{}, wg *sync.WaitGroup) {

	for {
		select {
		case err := <-errChan:
			log.Printf("error: %s", err)

		case <-shutdown:
			wg.Done()
			return

		default:
			// continue...
		}
	}
}
