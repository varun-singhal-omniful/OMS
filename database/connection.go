package database

import (
	"context"
	"fmt"
	"time"

	"github.com/omniful/go_commons/sqs"
	// "github.com/varun-singhal-omniful/oms-service/listeners"
	// "github.com/varun-singhal-omniful/oms-service/service"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Client
var Queue *sqs.Queue

func getDatabaseUri() string {
	return "mongodb://localhost:27017/OMS"
}

func ConnectMongo(c context.Context) {
	fmt.Println("Connecting to mongo...")
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	clientOptions := options.Client().ApplyURI(getDatabaseUri())
	var err error
	DB, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println("Error connecting to MongoDB:", err)
		return
	}
	err = DB.Ping(ctx, nil)
	if err != nil {
		fmt.Println("Failed to ping MongoDB:", err)
		return
	}

	fmt.Println("Successfully connected to MongoDB!")
}
func ConnectSqs(ctx context.Context) {
	acc := "124355654233"
	sqsConfig := sqs.GetSQSConfig(ctx, false, "ord", "eu-north-1", acc, "")
	fmt.Println(acc)
	url, err := sqs.GetUrl(ctx, sqsConfig, "MyQueue2")
	fmt.Println(*url)
	if err != nil {
		fmt.Println(err)
	}
	queueInstance, err := sqs.NewStandardQueue(ctx, "MyQueue2", sqsConfig)
	if err != nil {
		fmt.Println(err)
	}
	Queue = queueInstance
	fmt.Println(queueInstance)

}
