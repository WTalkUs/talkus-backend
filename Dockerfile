FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG FIREBASE_CREDENTIALS

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# Imagen final ligera
FROM alpine:3.21

WORKDIR /app

RUN echo "$FIREBASE_CREDENTIALS_B64" | base64 -d > /app/firebase.json

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]