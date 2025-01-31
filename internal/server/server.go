package server

import (
	"apollo-image-processor/internal/controller"
	"apollo-image-processor/internal/handler"
	"apollo-image-processor/internal/repository"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

type Server struct {
	db  *sql.DB
	api *http.Server
}

type Config struct {
	apihost    string
	apiport    string
	dbhost     string
	dbport     string
	dbuser     string
	dbpassword string
	dbname     string
}

func newConfig() Config {
	return Config{
		apihost:    os.Getenv("API_HOST"),
		apiport:    os.Getenv("API_PORT"),
		dbhost:     os.Getenv("DB_HOST"),
		dbport:     os.Getenv("DB_PORT"),
		dbuser:     os.Getenv("DB_USER"),
		dbpassword: os.Getenv("DB_PASSWORD"),
		dbname:     os.Getenv("DB_NAME"),
	}
}

func NewServer() (*Server, error) {
	server := Server{}
	config := newConfig()

	db, err := initDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	server.db = db

	api, err := initApiService(config, db)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize api server: %w", err)
	}
	server.api = api

	return &server, nil
}

func initApiService(config Config, db *sql.DB) (*http.Server, error) {

	router := chi.NewRouter()
	api := &http.Server{
		Addr:    ":" + config.apiport,
		Handler: router,
	}

	ir := repository.NewImageRepository(db)
	ic := controller.NewImageController(ir)
	ih := handler.NewImageHandler(ic)

	addRoutes(router, ih)

	log.Printf("listening on %s\n", api.Addr)
	if err := api.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return api, fmt.Errorf("error listening and serving: %w", err)
	}

	return api, nil
}

func initDB(config Config) (*sql.DB, error) {
	port, err := strconv.Atoi(config.dbport)
	if err != nil {
		return nil, fmt.Errorf("error getting port: %w", err)
	}

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.dbhost, port, config.dbuser, config.dbpassword, config.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, fmt.Errorf("could not connect to db: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to db: %w", err)
	}

	log.Println("connected to db...")

	return db, nil
}
