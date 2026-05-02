FROM golang:latest

WORKDIR /app

COPY go.mod .

RUN go mod download

COPY . .

RUN go build -o server ./cmd/server/main.go

CMD ["./server"]

