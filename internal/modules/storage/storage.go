package storage

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

type StorageService struct {
}

func NewService() *StorageService {
	return &StorageService{}
}

// folderPath does not include '/' as the last character
func (storage *StorageService) CreateFile(content []byte, fileName string, folderPath string) (string, string, error) {
	id := uuid.New().String()
	finalFileName := id + "_" + fileName
	filePath := folderPath + "/" + finalFileName
	file, err := os.Create(generateInternalStorageFilePath(filePath))

	if err != nil {
		return "", "", fmt.Errorf("could not create file: %w", err)
	}

	defer file.Close()

	_, err = file.Write(content)

	if err != nil {
		return "", "", fmt.Errorf("could not write to file: %w", err)
	}

	return finalFileName, generateExposedStorageFilePath(filePath), nil
}

func generateInternalStorageFilePath(filePath string) string {
	return fmt.Sprintf("%s%s", os.Getenv("INTERNAL_STORAGE_PREFIX"), filePath)
}

func generateExposedStorageFilePath(filePath string) string {
	return fmt.Sprintf("%s%s", os.Getenv("EXPOSED_STORAGE_PREFIX"), filePath)
}
