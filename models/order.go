package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type OrderItem struct {
	InventoryID primitive.ObjectID `bson:"inventory_id,omitempty" json:"inventory_id"`
	SKUID       string             `bson:"sku_id,omitempty" json:"sku_id"`
	Quantity    int                `bson:"quantity" json:"quantity"`
}

type Order struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id" `
	SellerID primitive.ObjectID `bson:"seller_id,omitempty" json:"seller_id"`
	HubID    primitive.ObjectID `bson:"hub_id,omitempty" json:"hub_id"`
	// CustomerID   primitive.ObjectID `bson:"customer_id,omitempty" json:"customer_id"`
	OrderNo      string             `bson:"order_no" json:"order_no"`
	CustomerName string             `bson:"customer_name" json:"customer_name"`
	OrderItems   []OrderItem        `bson:"order_items" json:"order_items" `
	Status       string             `bson:"status" json:"status"     `
	CreatedAt    primitive.DateTime `bson:"created_at,omitempty" json:"created_at"`
	UpdatedAt    primitive.DateTime `bson:"updated_at,omitempty" json:"updated_at" `
	DeletedAt    primitive.DateTime `bson:"deleted_at,omitempty" json:"deleted_at" `
}
