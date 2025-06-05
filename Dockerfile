FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Imagen final ligera
FROM alpine:3.21

WORKDIR /app

ENV FIREBASE_CREDENTIALS=""
RUN echo "$FIREBASE_CREDENTIALS" | base64 -d > /app/firebase.json

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]