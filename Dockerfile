FROM golang:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o server ./cmd/ordersystem

FROM scratch
COPY --from=builder /app/server .
COPY .env ./
COPY sql/migrations ./sql/migrations
EXPOSE 8080 8000 50051
ENTRYPOINT ["./server"]