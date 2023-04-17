FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY .env ./

RUN go mod download

COPY . .

WORKDIR /app/todo/cmd/server

RUN go build -o main .

EXPOSE 50051

CMD ["./main"]
