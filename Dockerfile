# Stage 1: Build
FROM golang:1.25-alpine AS builder
WORKDIR /app
# Copia i file dei moduli e scarica le dipendenze
COPY go.mod go.sum ./
RUN go mod download
# Copia il resto del codice
COPY . .
# Compila il binario statico
RUN CGO_ENABLED=0 GOOS=linux go build -o dispatcher ./cmd/server/main.go

# Stage 2: Immagine finale leggera
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/dispatcher .
EXPOSE 8080
CMD ["./dispatcher"]