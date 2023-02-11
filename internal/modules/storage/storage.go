package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type StorageService struct {
	cld *cloudinary.Cloudinary
}

func NewService() (*StorageService, error) {
	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUDINARY_CLOUD_NAME"), os.Getenv("CLOUDINARY_API_KEY"), os.Getenv("CLOUDINARY_API_SECRET"))

	if err != nil {
		return nil, err
	}

	return &StorageService{cld: cld}, nil
}

func (storage *StorageService) UploadFile(ctx context.Context, header *multipart.FileHeader) (string, string, error) {
	file, err := header.Open()
	defer file.Close()

	if err != nil {
		return "", "", err
	}

	res, err := storage.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(true)})

	if err != nil {
		return "", "", err
	}

	return res.AssetID, res.SecureURL, nil
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
