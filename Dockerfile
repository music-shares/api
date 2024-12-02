# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste du code
COPY . .

# Compiler l'application
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/api/main.go

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copier l'exécutable depuis le build stage
COPY --from=builder /app/main .

# Exposer le port de l'application
EXPOSE 8080

# Lancer l'application
CMD ["./main"]