# Tahap 1: Builder - Untuk kompilasi
FROM golang:1.23.2-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY .env .
# Kompilasi aplikasi Go menjadi binary statis
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./main.go

# Tahap 2: Final - Image akhir yang sangat kecil
FROM alpine:latest

WORKDIR /app
# Salin HANYA binary yang sudah dikompilasi dari tahap builder
COPY --from=builder /server .
COPY .env .

EXPOSE 4013
# Perintah untuk menjalankan aplikasi saat kontainer dimulai
CMD ["./server"]