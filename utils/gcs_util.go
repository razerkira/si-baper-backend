package utils

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/skip2/go-qrcode"
	"google.golang.org/api/option"
)

var storageClient *storage.Client

func InitGCS() {
	keyPath := os.Getenv("GCS_KEY_PATH")
	if keyPath == "" {
		log.Println("GCS Key Path tidak ditemukan, GCS mungkin tidak berfungsi.")
		return
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(keyPath))
	if err != nil {
		log.Fatalf("Gagal inisialisasi GCS Client: %v", err)
	}
	storageClient = client
	fmt.Println("GCS Client berhasil diinisialisasi!")
}

func GenerateAndUploadQRCode(itemCode string) (string, error) {
	bucketName := os.Getenv("GCS_BUCKET_NAME")
	if storageClient == nil || bucketName == "" {
		return "", fmt.Errorf("GCS tidak dikonfigurasi")
	}

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*30) 
	defer cancel()

	var qrPng []byte
	qrPng, err := qrcode.Encode(itemCode, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("gagal generate QR: %v", err)
	}

	fileName := fmt.Sprintf("qrcodes/qr_%s_%d.png", itemCode, time.Now().Unix())

	bucket := storageClient.Bucket(bucketName)
	object := bucket.Object(fileName)
	wc := object.NewWriter(ctx)
	wc.ContentType = "image/png"

	if _, err := wc.Write(qrPng); err != nil {
		return "", fmt.Errorf("gagal menulis data ke GCS: %v", err)
	}
	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("gagal menutup writer GCS: %v", err)
	}

	publicBaseURL := os.Getenv("GCS_PUBLIC_URL")
	qrCodeURL := fmt.Sprintf("%s/%s/%s", publicBaseURL, bucketName, fileName)

	return qrCodeURL, nil
}