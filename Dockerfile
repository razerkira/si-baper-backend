# Tahap 1: Builder (Menggunakan image Golang untuk meng-compile kode)
FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/api/main.go

# Tahap 2: Production (Membuat image yang super kecil dan ringan)
FROM alpine:latest

# Tambahkan tzdata (untuk zona waktu) DAN ca-certificates (untuk SSL/HTTPS)
RUN apk --no-cache add tzdata ca-certificates

WORKDIR /app

COPY --from=builder /app/server .

EXPOSE 8080

CMD ["./server"]