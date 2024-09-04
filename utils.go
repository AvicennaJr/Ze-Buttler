package main

import (
	"io"
	"log"
	"os"
)

func imageToBytes(imagePath string) ([]byte, error) {

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatal("Error opening image file:", err)
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return imageBytes, nil

}
