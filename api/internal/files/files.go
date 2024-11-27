package files

import (
	"bytes"
	"image"
	"image/jpeg"
	"log/slog"
	"os"
	"time"

	_ "image/png" // Поддержка PNG
	"math/rand"

	"github.com/nfnt/resize"
)

type FileManager struct {
}

func New() *FileManager {
	return &FileManager{}
}

func generateRandomString() string {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := make([]byte, 10)
	for i := range b {
		b[i] = letters[randomizer.Intn(len(letters))]
	}
	return string(b)
}

func compressImage(photo []byte) ([]byte, error) {
	// Декодируем изображение из байтов
	img, _, err := image.Decode(bytes.NewReader(photo))
	if err != nil {
		return nil, err
	}

	// Изменяем размер изображения (сохраняем пропорции, уменьшаем в 2 раза)
	bounds := img.Bounds()
	newWidth := uint(bounds.Dx() / 2)
	newHeight := uint(bounds.Dy() / 2)

	// Сжимаем изображение
	newImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// Кодируем изображение в JPEG с заданным качеством
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImg, &jpeg.Options{Quality: 90}) // Установите качество 85 (или любое другое)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (f *FileManager) UploadImage(file []byte, name string) (string, error) {
	file, err := compressImage(file)
	if err != nil {
		slog.Error("Error compressing image: " + err.Error())
		return "", err
	}

	filePath := os.Getenv("FILE_PATH") + "/images/" + name + generateRandomString() + ".png"
	newFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}

	defer newFile.Close()

	_, err = newFile.Write(file)
	if err != nil {
		return "", err
	}

	return filePath, nil

}

func (f *FileManager) SaveImage(file []byte, name string) error {
	filePath := os.Getenv("FILE_PATH") + "/images/" + name + ".jpg"
	newFile, err := os.Create(filePath)
	if err != nil {
		slog.Error("Error creating file: " + err.Error())
		return err
	}

	defer newFile.Close()

	_, err = newFile.Write(file)
	if err != nil {
		slog.Error("Error writing file: " + err.Error())
		return err
	}

	return nil

}
