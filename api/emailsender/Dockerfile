# Stage 1: Build the application
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o emailsender ./cmd/emailsender

# Stage 2: Create a minimal runtime image
FROM scratch

WORKDIR /app

COPY --from=builder /app/emailsender .

CMD ["./emailsender"]