package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/tatrasoft/fyp-backend/database"
	"github.com/tatrasoft/fyp-backend/models"
	"github.com/tatrasoft/fyp-backend/proto"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	collectionItems = "items"
)

type ItemsServerService struct{}

var itemsColHelper database.CollectionHelper

func SetColHelper(helper database.DBHelper) {
	itemsColHelper = helper.Collection(collectionItems)
}

// CreateItem creates new item entry in the db
func (iss *ItemsServerService) CreateItem(
	ctx context.Context,
	req *proto.CreateItemReq) (*proto.CreateItemRes, error) {
	// to access the struct with nil check
	item := req.GetItem()
	// converting Item type to BSON
	data := models.Item{
		ID:                  primitive.ObjectID{},
		ItemName:            item.GetItemName(),
		ItemDescription:     item.GetItemDescription(),
		ItemImg:             item.GetItemImg(),
		ItemPrice:           item.GetItemPrice(),
		EstimatePrepareTime: item.GetEstimatePrepareTime(),
		CreatedAt:           time.Now().Unix(),
		UpdatedAt:           time.Now().Unix(),
		DeletedAt:           0,
	}
	// insert data into database, result contains newly generated Object ID for the new document
	result, err := itemsColHelper.InsertOne(ctx, &data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	// convert object id to its string counterpart
	item.Id = result.(primitive.ObjectID).Hex()

	return &proto.CreateItemRes{Item: item}, nil
}

// ReadItem get required item from the db
func (iss *ItemsServerService) ReadItem(ctx context.Context, req *proto.ReadItemReq) (*proto.ReadItemRes, error) {
	// convert string id from proto to mongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("could not convert to objectID: %v", err))
	}
	result := itemsColHelper.FindOne(ctx, bson.M{"_id": oid})

	// create empty Item object to write our decode result into
	data := models.Item{}
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("could not find item with Object Id %s: %v", req.GetId(), err))
	}

	response := &proto.ReadItemRes{
		Item: &proto.Item{
			Id:                   oid.Hex(),
			ItemName:             data.ItemName,
			ItemDescription:      data.ItemDescription,
			ItemImg:              data.ItemImg,
			ItemPrice:            data.ItemPrice,
			EstimatePrepareTime:  data.EstimatePrepareTime,
			CreateAt:				data.CreatedAt,
			UpdatedAt: data.UpdatedAt,
			DeletedAt: data.DeletedAt,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
	}

	return response, nil
}

// DeleteItem removes required entry from the database
func (iss *ItemsServerService) DeleteItem(ctx context.Context, req *proto.DeleteItemReq) (*proto.DeleteItemRes, error) {
	// Get the ID (string) from the request message and convert it to an Object ID
	oid, err := primitive.ObjectIDFromHex(req.GetId())
	// Check for errors
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Could not convert to ObjectId: %v", err))
	}
	// DeleteOne returns DeleteResult which is a struct containing the amount of deleted docs (in this case only 1 always)
	// So we return a boolean instead
	_, err = itemsColHelper.DeleteOne(ctx, bson.M{"_id": oid})
	// Check for errors
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete item with id %s: %v", req.GetId(), err))
	}

	// Return response with success: true if no error is thrown (and thus document is removed)
	return &proto.DeleteItemRes{
		Success: true,
	}, nil
}

// updating entry in the db
func (iss *ItemsServerService) UpdateItem(ctx context.Context, req *proto.UpdateItemReq) (*proto.UpdateItemRes, error) {
	// Get the item data from the request
	item := req.GetItem()

	// Convert the Id string to a MongoDB ObjectId
	oid, err := primitive.ObjectIDFromHex(item.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Could not convert the supplied item id to a MongoDB ObjectId: %v", err),
		)
	}

	// Convert the data to be updated into an unordered Bson document
	update := bson.M{
		"item_name":        item.GetItemName(),
		"item_description": item.GetItemDescription(),
		"item_img":         item.GetItemImg(),
		"item_price":       item.GetItemPrice(),
		"estimate_prepare_time": item.GetEstimatePrepareTime(),
		"create_at": item.GetCreateAt(),
		"updated_at": item.GetUpdatedAt(),
		"deleted_at": item.GetDeletedAt(),
	}

	// Convert the oid into an unordered bson document to search by id
	filter := bson.M{"_id": oid}

	// Result is the BSON encoded result
	result := itemsColHelper.FindOneAndUpdate(ctx, filter, bson.M{"$set": update})

	// Decode result and write it to 'decoded'
	decoded := models.Item{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find item with supplied ID: %v", err),
		)
	}

	return &proto.UpdateItemRes{
		Item: &proto.Item{
			Id:                   decoded.ID.Hex(),
			ItemName:             decoded.ItemName,
			ItemDescription:      decoded.ItemDescription,
			ItemImg:              decoded.ItemImg,
			ItemPrice:            decoded.ItemPrice,
			EstimatePrepareTime: decoded.EstimatePrepareTime,
			CreateAt: decoded.CreatedAt,
			UpdatedAt: time.Now().Unix(),
			DeletedAt: decoded.DeletedAt,
			XXX_NoUnkeyedLiteral: struct{}{},
			XXX_unrecognized:     nil,
			XXX_sizecache:        0,
		},
	}, nil
}

// list all entries form the db
func (iss *ItemsServerService) ListItems(req *proto.ListItemReq, stream proto.ItemService_ListItemsServer) error {
	// Initiate a Item type to write decoded data to
	data := &models.Item{}
	// collection.Find returns a cursor for our (empty) query
	cursor, err := itemsColHelper.List(context.Background(), bson.M{})
	if err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	// An expression with defer will be called at the end of the function
	defer cursor.Close(context.Background())
	// cursor.Next() returns a boolean, if false there are no more items and loop will break
	for cursor.Next(context.Background()) {
		// Decode the data at the current pointer and write it to data
		err := cursor.Decode(data)
		// check error
		if err != nil {
			return status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		// If no error is found send item over stream
		err = stream.Send(&proto.ListItemRes{
			Item: &proto.Item{
				Id:                   data.ID.Hex(),
				ItemName:             data.ItemName,
				ItemDescription:      data.ItemDescription,
				ItemImg:              data.ItemImg,
				ItemPrice:            data.ItemPrice,
				XXX_NoUnkeyedLiteral: struct{}{},
				XXX_unrecognized:     nil,
				XXX_sizecache:        0,
			},
		})
		if err != nil {
			return status.Errorf(codes.Internal, fmt.Sprintf("unable to send item over the stream: %v", err))
		}
	}
	// Check if the cursor has any errors
	if err := cursor.Err(); err != nil {
		return status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}

	return nil
}
