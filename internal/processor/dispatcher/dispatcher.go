package dispatcher

import (
	"apollo-image-processor/internal/models"
	procrepository "apollo-image-processor/internal/processor/repository"
	"apollo-image-processor/internal/processor/worker"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/streadway/amqp"
)

type Dispatcher struct {
	db             *sql.DB
	rmqpool        *sync.Pool
	workerPoolSize int
	prefetchCount  int
	shutdown       chan struct{}
	wg             sync.WaitGroup
	pr             procrepository.ProcessorRepository
}

func NewDispatcher(rmqpool *sync.Pool, pr procrepository.ProcessorRepository) (*Dispatcher, error) {
	workerPoolSize, err := strconv.Atoi(os.Getenv("WORKER_POOL_SIZE"))
	if err != nil {
		return nil, fmt.Errorf("err setting worker pool size: %w", err)
	}

	prefetchCount, err := strconv.Atoi(os.Getenv("RMQ_PREFETCH_COUNT"))
	if err != nil {
		return nil, fmt.Errorf("error setting prefetch count: %w", err)
	}

	shutdown := make(chan struct{})

	return &Dispatcher{
		rmqpool:        rmqpool,
		workerPoolSize: workerPoolSize,
		prefetchCount:  prefetchCount,
		shutdown:       shutdown,
		wg:             sync.WaitGroup{},
		pr:             pr,
	}, nil
}

func (d *Dispatcher) Start() error {

	rmqChan := d.rmqpool.Get().(*amqp.Channel)
	err := rmqChan.Qos(
		d.prefetchCount, // prefetchCount
		0,               // size
		false,           // global
	)
	if err != nil {
		return fmt.Errorf("error setting QoS: %w", err)
	}

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

	go func() {
		for msg := range msgs {

			jobsChan <- msg

		}
	}()

	// start up our workers
	for w := 0; w < d.workerPoolSize; w++ {
		d.wg.Add(1)
		go worker.ImgWorker(jobsChan, resChan, errChan, d.shutdown, &d.wg, d.pr)
	}

	d.wg.Add(1)
	go worker.ResWorker(resChan, d.shutdown, rmqChan, &d.wg)

	d.wg.Add(1)
	go worker.ErrWorker(errChan, d.shutdown, &d.wg)

	log.Println("Waiting for messages...")

	return nil
}

func (d *Dispatcher) Shutdown(ctx context.Context) error {

	log.Printf("shutting down the dispatcher...")

	close(d.shutdown)

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All workers shut down successfully")
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())

	}

}
