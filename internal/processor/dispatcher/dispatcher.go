package dispatcher

import (
	"apollo-image-processor/internal/models"
	"apollo-image-processor/internal/processor/worker"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
)

func InitDispatcher(rmqpool *sync.Pool, db *sql.DB) error {

	workerPoolSize, err := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
	if err != nil {
		return fmt.Errorf("err setting worker pool size: %w", err)
	}

	prefetchCount, err := strconv.Atoi(os.Getenv("RMQ_PREFETCH_COUNT"))
	if err != nil {
		return fmt.Errorf("error setting prefetch count: %w", err)
	}

	rmqChan := rmqpool.Get().(*amqp.Channel)
	err = rmqChan.Qos(
		prefetchCount, // prefetchCount
		0,             // size
		false,         // global
	)
	if err != nil {
		return fmt.Errorf("error setting QoS: %w", err)
	}
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
		return fmt.Errorf("error declaring queue: %w", err)
	}
	log.Printf("queue status: %+v", queue)

	msgs, err := rmqChan.Consume(
		models.QueueName, // queue
		"",               // consumer
		false,            // auto ack
		false,            // exclusive
		false,            // no local
		false,            // no wait
		nil,              //args
	)
	if err != nil {
		return fmt.Errorf("error getting messages: %w", err)
	}

	// create channel that will hold our messages
	jobsChan := make(chan amqp.Delivery, 10)
	resChan := make(chan amqp.Delivery)
	errChan := make(chan error)

	forever := make(chan bool)
	go func() {
		for msg := range msgs {

			jobsChan <- msg

		}
	}()

	// start up our workers
	for w := 0; w < workerPoolSize; w++ {
		go worker.ImgWorker(jobsChan, resChan, errChan, db)
	}

	go worker.ResWorker(resChan, rmqChan)
	go worker.ErrWorker(errChan)

	log.Println("Waiting for messages...")
	<-forever

	return nil
}
