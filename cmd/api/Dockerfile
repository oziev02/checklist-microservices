FROM golang:1.21 AS builder
WORKDIR /app
COPY ../../go.mod ../../go.sum ./
RUN go mod download
COPY ../.. .
RUN go build -o api ./cmd/api
FROM golang:1.21
WORKDIR /app
COPY --from=builder /app/api .
CMD ["./api"]