# Build stage
FROM --platform=linux/amd64 golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Spécifier explicitement l'architecture amd64
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o main cmd/api/main.go

# Final stage
FROM --platform=linux/amd64 alpine:3.18

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 10000
CMD ["./main"]