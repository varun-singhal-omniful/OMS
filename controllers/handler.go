package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/http"
	interservice_client "github.com/omniful/go_commons/interservice-client"
	"github.com/varun-singhal-omniful/oms-service/database"
	"github.com/varun-singhal-omniful/oms-service/kafka"
	"github.com/varun-singhal-omniful/oms-service/models"
	"github.com/varun-singhal-omniful/oms-service/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderRequest struct {
	FilePath string             `json:"file_path"`
	SellerID primitive.ObjectID `json:"seller_id"`
	HubID    primitive.ObjectID `json:"hub_id"`
}

func BulkOrders(c *gin.Context) {
	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	if _, err := os.Stat(req.FilePath); os.IsNotExist(err) {
		c.JSON(400, gin.H{"error": "File not found"})
		return
	}
	service.SetProducer(c, database.Queue, req.FilePath)

	orders, err := performcsvopr(req.FilePath, req.SellerID, req.HubID)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, order := range orders {
		if err := storeOrder(order); err != nil {
			c.JSON(500, gin.H{"error": "Failed to save order"})
			return
		}
	}

	c.JSON(200, gin.H{"message": "Orders uploaded successfully", "total_orders": len(orders)})
}

func performcsvopr(filePath string, sellerID, hubID primitive.ObjectID) ([]*models.Order, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer file.Close()
	orderGroups := make(map[string]*models.Order)
	Csv, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(filePath),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}
	err = Csv.InitializeReader(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize CSV reader: %v", err)
	}
	for !Csv.IsEOF() {
		var records csv.Records
		records, err := Csv.ReadNextBatch()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Processing records:")
		fmt.Println(records)
		for _, record := range records {
			orderNo := record[0]
			customerName := record[1]
			skuID := record[2]
			quantityStr := record[3]
			quantity, err := strconv.Atoi(quantityStr)
			if err != nil {
				return nil, fmt.Errorf("invalid quantity %s: %v", quantityStr, err)
			}
			orderKey := fmt.Sprintf("%s-%s", orderNo, customerName)
			order, exists := orderGroups[orderKey]
			if !exists {

				now := primitive.NewDateTimeFromTime(time.Now())
				order = &models.Order{
					ID:           primitive.NewObjectID(),
					SellerID:     sellerID,
					HubID:        hubID,
					CustomerName: customerName,
					OrderNo:      orderNo,
					OrderItems:   []models.OrderItem{},
					Status:       "on_hold",
					CreatedAt:    now,
					UpdatedAt:    now,
				}
				orderGroups[orderKey] = order
			}
			orderItem := models.OrderItem{
				SKUID:    skuID,
				Quantity: quantity,
			}
			order.OrderItems = append(order.OrderItems, orderItem)
		}
	}
	var orders []*models.Order
	for _, order := range orderGroups {
		orders = append(orders, order)
	}

	fmt.Println("Final orders:")
	for _, order := range orders {
		orderData, _ := json.Marshal(order)
		kafka.PublishMessageToKafka(orderData, order.ID.Hex())
		config := interservice_client.Config{
			ServiceName: "user-service",
			BaseURL:     "http://localhost:8081/api/V1/sku/",
			Timeout:     5 * time.Second,
		}
		client, _ := interservice_client.NewClientWithConfig(config)
		url := config.BaseURL + "validate"
		body := map[string]string{
			"hub_id": "",
			"skus":   "",
		}
		bodyBytes, _ := json.Marshal(body)

		// if err != nil {
		// 	return false
		// }
		req := &http.Request{
			Url:     url, // Use configured URL
			Body:    bytes.NewReader(bodyBytes),
			Timeout: 7 * time.Second,
			Headers: map[string][]string{
				"Content-Type": {"application/json"},
			},
		}

		// if err != nil {
		// 	return
		// }

		resp, _ := client.Get(req, "/")
		fmt.Println("resp is", resp)
		fmt.Printf("Order No: %s, Customer: %s, Total Items: %d\n", order.OrderNo, order.CustomerName, len(order.OrderItems))
	}

	return orders, nil
}

// storeOrder inserts a single order into MongoDB
func storeOrder(order *models.Order) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.DB.Database("OMS").Collection("orders")

	_, err := collection.InsertOne(ctx, order)
	return err
}
