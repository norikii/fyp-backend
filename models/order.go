package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Order struct {
	ID primitive.ObjectID `bson:"_id"`
	TableID int32 `bson:"table_id"`
	UserID primitive.ObjectID `bson:"user_id"`
	Items []Item `bson:"items"`
	TotalPrice float32 `bson:"total_price"`
	PayedAt int64 `bson:"payed_at"`
	IsReadyAt int64 `bson:"is_ready_at"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
	DeletedAt int64 `bson:"deleted_at"`
}
