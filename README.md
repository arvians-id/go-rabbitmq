# Example of RabbitMQ in Golang

This is an example of implementation of RabbitMQ in Golang.

## Requirements
Make sure you have installed all of the following prerequisites on your development machine:
* [Git](https://git-scm.com/)
* [Go](https://golang.org/)
* [Docker](https://www.docker.com/)
* [Golang Migrate](https://github.com/golang-migrate/migrate)
* [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial)
* [Air](https://github.com/cosmtrek/air)

## Installation
This project requires [Go](https://golang.org/) v1.20+ to run.

```bash
# Clone this project
$ git clone https://github.com/arvians-id/go-rabbitmq.git

# Move to project directory
$ cd go-rabbitmq

# Copy .env.example to .env
$ cp category/.env.example category/.env \
    && cp todo/.env.example todo/.env \
    && cp user/.env.example user/.env \
    && cp worker/.env.example worker/.env \
    && cp category_todo/.env.example category_todo/.env \
    && cp gateway/.env.example gateway/.env

# Move to each directory and install the dependencies
$ go mod download
# or
$ go mod tidy

# Migrate tables
$ make migrate table=go_rabbitmq verbose=up
$ make migrate table=go_rabbitmq_test verbose=up

# After installing protocol buffers, generate the proto file
$ make pb
```


## Run Application

### On Development
```bash
# Run all containers
$ docker-compose -f docker-compose.dev.yml up -d

# Run all services
# Warning: Make sure you have running all services in different terminal
$ go run cmd/category/main.go
$ go run cmd/todo/main.go
$ go run cmd/user/main.go
$ go run cmd/worker/main.go
$ go run cmd/category_todo/main.go
$ go run cmd/gateway/main.go

# or you can run all services in different terminal with air
$ air

# graphql playground, open in browser
http://localhost:3000/playground

# rest api, open in postman or curl
$ curl -X GET http://localhost:3000/api/users

# Run test
# Warning: Make sure you have run on development mode
$ make test
```

### On Production
```bash
# Build all services
$ make build

# Run database on docker
$ docker-compose up -d
```

### Technologies Used
* [Go](https://golang.org/)
* [Docker](https://www.docker.com/)
* [gRPC](https://grpc.io/)
* [GraphQL](https://graphql.org/)
* [RabbitMQ](https://www.rabbitmq.com/)
* [PostgreSQL](https://www.postgresql.org/)
* [Redis](https://redis.io/)
* [Open Telementry (Jaeger)](https://opentelemetry.io/)