package Init

import (
	"context"
	"fmt"

	k "github.com/omniful/go_commons/kafka"
	"github.com/varun-singhal-omniful/oms-service/database"
	"github.com/varun-singhal-omniful/oms-service/kafka"
)

func InitializeDB(c context.Context) {
	database.ConnectMongo(c)
}

func InitializeSqs(c context.Context) {
	database.ConnectSqs(c)
}
func InitializeKafkaProducer(ctx context.Context) {
	kafkaBrokers := make([]string, 1)
	kafkaBrokers[0] = "localhost:9092"
	kafkaClientID := "tenant-service"
	kafkaVersion := "2.0.0"
	fmt.Print("kafka version is : ", kafkaVersion, "\n")

	producer := k.NewProducer(
		k.WithBrokers(kafkaBrokers),
		k.WithClientID(kafkaClientID),
		k.WithKafkaVersion(kafkaVersion),
	)
	fmt.Println("Initialized Kafka Producer")
	kafka.SetProducer(producer)

}
