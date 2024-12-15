package services

import (
	"log"

	"github.com/streadway/amqp"
)

var rabbitConn *amqp.Connection
var rabbitChannel *amqp.Channel

func InitRabbitMQ() {
	var err error
	rabbitConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	rabbitChannel, err = rabbitConn.Channel()
	if err != nil {
		log.Fatal("Failed to open RabbitMQ channel:", err)
	}

	_, err = rabbitChannel.QueueDeclare(
		"image_queue", true, false, false, false, nil,
	)
	if err != nil {
		log.Fatal("Failed to declare RabbitMQ queue:", err)
	}
}

func QueueProductImages(images []string) {
	for _, image := range images {
		err := rabbitChannel.Publish(
			"", "image_queue", false, false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(image),
			},
		)
		if err != nil {
			log.Println("Failed to publish image URL:", image)
		}
	}
}
