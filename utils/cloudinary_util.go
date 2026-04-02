package utils

import (
	"bytes"
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

// UploadQRCodeToCloudinary mengunggah file gambar ke Cloudinary dan mengembalikan URL publiknya
func UploadQRCodeToCloudinary(fileName string, fileData []byte) (string, error) {
	// Koneksi ke Cloudinary menggunakan Environment Variables
	cld, err := cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
	if err != nil {
		return "", err
	}

	ctx := context.Background()

	// Unggah file berupa byte array
	resp, err := cld.Upload.Upload(ctx, bytes.NewReader(fileData), uploader.UploadParams{
		Folder:   "sibaper-qrcodes", // Otomatis membuat folder di Cloudinary
		PublicID: fileName,          // Nama file (tanpa ekstensi .png)
	})
	if err != nil {
		return "", err
	}

	// Kembalikan URL HTTPS yang aman
	return resp.SecureURL, nil
}