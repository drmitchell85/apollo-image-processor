package apiserver

import (
	"apollo-image-processor/internal/api/controller"
	"apollo-image-processor/internal/api/handler"
	"apollo-image-processor/internal/api/messenger"
	"apollo-image-processor/internal/api/repository"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

type APIServer struct {
	db      *sql.DB
	api     *http.Server
	amqp    *amqp.Connection
	rmqpool *sync.Pool
}

type Config struct {
	APIhost    string
	APIport    string
	DBhost     string
	DBport     string
	DBuser     string
	DBpassword string
	DBname     string
	RMQprefix  string
	RMQuser    string
	RMQhost    string
	RMQport    string
}

func newConfig() Config {
	return Config{
		APIhost:    os.Getenv("API_HOST"),
		APIport:    os.Getenv("API_PORT"),
		DBhost:     os.Getenv("DB_HOST"),
		DBport:     os.Getenv("DB_PORT"),
		DBuser:     os.Getenv("DB_USER"),
		DBpassword: os.Getenv("DB_PASSWORD"),
		DBname:     os.Getenv("DB_NAME"),
		RMQprefix:  os.Getenv("RMQ_PREFIX"),
		RMQuser:    os.Getenv("RMQ_USER"),
		RMQhost:    os.Getenv("RMQ_HOST"),
		RMQport:    os.Getenv("RMQ_PORT"),
	}
}

func (s *APIServer) Start() error {
	log.Printf("listening on %s\n", s.api.Addr)
	if err := s.api.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error listening and serving: %w", err)
	}

	return nil
}

func (s *APIServer) Shutdown(ctx context.Context) {

	if err := s.api.Shutdown(ctx); err != nil {
		log.Printf("error shutting down server: %v", err)
	}

	if err := s.db.Close(); err != nil {
		log.Printf("error closing db connection: %v", err)
	}

	if err := s.amqp.Close(); err != nil {
		log.Printf("error closing amqp connection: %v", err)
	}

}

func NewServer() (*APIServer, error) {
	server := APIServer{}
	config := newConfig()

	db, err := initDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	server.db = db

	amqp, rmqpool, err := initRMQ(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize rabbitMQ: %w", err)
	}
	server.amqp = amqp
	server.rmqpool = rmqpool

	api, err := initApiService(config, db, rmqpool)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize api server: %w", err)
	}
	server.api = api

	return &server, nil
}

func initApiService(config Config, db *sql.DB, rmqpool *sync.Pool) (*http.Server, error) {

	router := chi.NewRouter()
	api := &http.Server{
		Addr:    ":" + config.APIport,
		Handler: router,
	}

	mq := messenger.NewMessengerQueue(rmqpool)
	ir := repository.NewImageRepository(db)
	ic := controller.NewImageController(ir, mq)
	ih := handler.NewImageHandler(ic)

	addRoutes(router, ih)

	return api, nil
}

func initDB(config Config) (*sql.DB, error) {
	port, err := strconv.Atoi(config.DBport)
	if err != nil {
		return nil, fmt.Errorf("error getting port: %w", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.DBhost, port, config.DBuser, config.DBpassword, config.DBname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	log.Println("connected to db")

	return db, nil
}

func initRMQ(config Config) (*amqp.Connection, *sync.Pool, error) {

	// "amqp://guest:guest@localhost:port/"
	rmqUrl := fmt.Sprintf(
		"%s%s:%s@%s:%s/",
		config.RMQprefix,
		config.RMQuser,
		config.RMQuser,
		config.RMQhost,
		config.RMQport,
	)

	amqpConfig := amqp.Config{
		Heartbeat: time.Second * 60,
	}

	conn, err := amqp.DialConfig(rmqUrl, amqpConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("error connecting to rabbitMQ: %w", err)
	}
	log.Println("connected to rabbitMQ")

	var rmqPool = &sync.Pool{
		New: func() interface{} {
			channel, err := conn.Channel()
			if err != nil {
				log.Printf("error creating rabbitMQ channel: %s", err)
			} else {
				log.Printf("rabbitMQ channel created")
			}
			return channel
		},
	}
	log.Println("rabbitMQ pool created")

	return conn, rmqPool, nil
}
