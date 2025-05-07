FROM golang:1.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd

FROM debian:bookworm-slim

WORKDIR /app
COPY --from=builder /app/main .

COPY wait-for-it.sh /wait-fot-it.sh
RUN chmod +x /wait-fot-it.sh

EXPOSE 8080

CMD ["./main"]