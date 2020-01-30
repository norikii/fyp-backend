package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	TableID int32 `bson:"table_id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Items []Item `bson:"items"`
	TotalPrice float32 `bson:"total_price"`
	IsPayed bool `bson:"is_payed"`
	IsReady bool `bson:"is_ready"`
}
