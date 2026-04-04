package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func SaveQRCodeLocally(itemCode string, fileData []byte) (string, error) {
	fileName := fmt.Sprintf("%s.png", itemCode)
	uploadDir := "./uploads/qrcodes"
	filePath := filepath.Join(uploadDir, fileName)

	err := os.WriteFile(filePath, fileData, 0644)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/qrcodes/%s", fileName), nil
}