package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	ItemName string `bson:"table_id"`
	ItemDescription string `bson:"item_description"`
	ItemImg string `bson:"item_img"`
	ItemPrice float32 `bson:"item_price"`
}
