package service

import (
	"context"
	"fmt"
	"log"

	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go/aws"

	// "github.com/joho/godotenv"
	"github.com/omniful/go_commons/sqs"
)

var NewProducer = &sqs.Publisher{}

func SetProducer(ctx context.Context, queue *sqs.Queue, message string) {
	fmt.Println("message is", message)
	NewProducer = sqs.NewPublisher(queue)
	newmessage := &sqs.Message{
		GroupId:       "group-123",
		Value:         []byte(message),
		ReceiptHandle: "receipt-abc",
		Attributes:    map[string]string{"key1": "value1", "key2": "value2"},
	}
	err := NewProducer.Publish(ctx, newmessage)
	if err != nil {
		log.Fatal("Error in publishing the message", err)
	}

}
