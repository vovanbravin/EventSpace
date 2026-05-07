FROM golang:latest

WORKDIR /app

COPY ../../go.mod .

RUN go mod download

COPY ../.. .

RUN go build -o api-gateway ./cmd/api-gateway/main.go

CMD ["./api-gateway"]

