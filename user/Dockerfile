FROM golang:1.20-alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod tidy

COPY . .

EXPOSE 50050

CMD ["go", "run", "cmd/server/main.go"]