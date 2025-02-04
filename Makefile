include .env
export

.PHONY: run-api
run-api:
	go run cmd/api/main.go

.PHONY: run-processor
run-processor:
	go run cmd/processor/main.go

.PHONY: docker-build
docker-build:
	docker build -t apollo-image-processor:multistage -f Dockerfile.multistage .

.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 --env-file=.env apollo-image-processor:multistage