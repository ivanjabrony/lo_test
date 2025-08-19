FROM golang:1.24.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o task_lo ./cmd/main.go

FROM alpine:3.18

WORKDIR /app

COPY --from=builder /app/task_lo .

EXPOSE 8080

ENTRYPOINT ["./task_lo"]