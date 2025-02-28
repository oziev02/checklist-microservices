FROM golang:1.21 AS builder
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY ../.. .
RUN go build -o kafka ./cmd/kafka
FROM golang:1.21
WORKDIR /app
COPY --from=builder /app/kafka .
CMD ["./kafka"]