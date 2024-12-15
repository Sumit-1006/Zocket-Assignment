package services

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"

	"github.com/nfnt/resize"
)

func ProcessImages() {
	msgs, err := rabbitChannel.Consume("image_queue", "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to consume RabbitMQ messages:", err)
	}

	for msg := range msgs {
		imageURL := string(msg.Body)
		resp, err := http.Get(imageURL)
		if err != nil {
			log.Println("Failed to download image:", imageURL)
			continue
		}
		defer resp.Body.Close()

		img, _, err := image.Decode(resp.Body)
		if err != nil {
			log.Println("Failed to decode image:", imageURL)
			continue
		}

		resizedImage := resize.Resize(200, 0, img, resize.Lanczos3)
		var buf bytes.Buffer
		err = jpeg.Encode(&buf, resizedImage, nil)
		if err != nil {
			log.Println("Failed to compress image:", imageURL)
			continue
		}

		// Save compressed image to a file (or S3)
		file, err := os.Create("compressed_" + imageURL)
		if err != nil {
			log.Println("Failed to save compressed image:", imageURL)
			continue
		}
		defer file.Close()

		buf.WriteTo(file)
	}
}
