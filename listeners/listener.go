package listeners

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/joho/godotenv"
)

func SetConsumer() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	queue_url := os.Getenv("QUEUE_URL")
	if queue_url == "" {
		log.Fatal("No QUEUE_URL specified")
	}
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	sqsClient := sqs.NewFromConfig(cfg)
	for {
		recieveMessages := &sqs.ReceiveMessageInput{
			QueueUrl:            aws.String(queue_url),
			MaxNumberOfMessages: 10,
			WaitTimeSeconds:     5,
		}
		resp, err := sqsClient.ReceiveMessage(ctx, recieveMessages)
		if err != nil {
			log.Fatal(err)
		}
		for _, message := range resp.Messages {
			fmt.Println("received message is : ", *message.Body)

			_, err := sqsClient.DeleteMessage(ctx, &sqs.DeleteMessageInput{
				QueueUrl:      aws.String(queue_url),
				ReceiptHandle: message.ReceiptHandle,
			})

			if err != nil {
				log.Fatal(err)
			} else {
				fmt.Println("Message deleted from queue")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
