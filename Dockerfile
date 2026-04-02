FROM golang:1.25.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .


RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/api/main.go

FROM alpine:latest

RUN apk --no-cache add tzdata

WORKDIR /app

COPY --from=builder /app/server .

COPY gcs-key.json .

COPY .env .

EXPOSE 8080

CMD ["./server"]