# tell Docker what base image to use
FROM golang:1.23

# create a directory inside the image that we are building
WORKDIR /app

# copy our mod and sum files
COPY go.mod go.sum ./

# install modules in the image
RUN go mod download

# copy source code into the image
COPY . .

# compile the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /apollo-image-processor ./cmd/api

# expose a port for our application
EXPOSE 8080

# tell docker what to run when image starts in the container
CMD ["/apollo-image-processor"]