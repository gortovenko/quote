FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ecs_server ./ecs/ecs_main.go

CMD ["./ecs_server"]