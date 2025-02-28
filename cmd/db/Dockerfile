FROM golang:1.21 AS builder
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY ../.. .
RUN go build -o db ./cmd/db
FROM golang:1.21
WORKDIR /app
COPY --from=builder /app/db .
CMD ["./db"]