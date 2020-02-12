package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID primitive.ObjectID `bson:"_id"`
	ItemName string `bson:"table_id"`
	ItemDescription string `bson:"item_description"`
	ItemImg string `bson:"item_img"`
	ItemPrice float32 `bson:"item_price"`
	EstimatePrepareTime  int64 `bson:"estimate_prepare_time"`
	CreatedAt int64 `bson:"created_at"`
	UpdatedAt int64 `bson:"updated_at"`
	DeletedAt int64 `bson:"deleted_at"`
}
