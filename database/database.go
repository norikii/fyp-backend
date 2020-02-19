package database

import (
	"context"

	"github.com/tatrasoft/fyp-backend/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

type CollectionHelper interface {
	FindOne(context.Context, interface{}) SingleResultHelper
	InsertOne(context.Context, interface{}) (interface{}, error)
	DeleteOne(ctx context.Context, filer interface{}) (int64, error)
	FindOneAndUpdate(context.Context, interface{}, interface{}) SingleResultHelper
	List(context.Context, interface{}) (*mongo.Cursor, error)
}

type SingleResultHelper interface {
	Decode(v interface{}) error
}

type ClientHelper interface {
	Database(string) DBHelper
	Connect(ctx context.Context) error
	CloseConnection(ctx context.Context) error
	StartSession() (mongo.Session, error)
	Ping(ctx context.Context) error
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDB struct {
	db *mongo.Database
}

type mongoCollection struct {
	collection *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoSession struct {
	mongo.Session
}

//func NewClient(cnf *config.DBConfig) (ClientHelper, error) {
//	c, err := mongo.NewClient(options.Client().SetAuth(
//		options.Credential{
//			AuthSource:              cnf.DatabaseName,
//			Username:                cnf.Username,
//			Password:                cnf.Password,
//		}).ApplyURI(cnf.Url))
//
//	return &mongoClient{cl: c}, err
//}

func NewClient(cnf *config.DBConfig) (ClientHelper, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(cnf.Url))

	return &mongoClient{cl: c}, err
}

func NewDatabase(cnf *config.DBConfig, client ClientHelper) DBHelper {
	return client.Database(cnf.DatabaseName)
}

func (mc *mongoClient) Database(dbName string) DBHelper {
	db := mc.cl.Database(dbName)

	return &mongoDB{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()

	return &mongoSession{session}, err
}

func (mc *mongoClient) Ping(ctx context.Context) error {
	return mc.cl.Ping(ctx, nil)
}

func (mc *mongoClient) Connect(ctx context.Context) error {
	return mc.cl.Connect(ctx)
}

func (mc *mongoClient) CloseConnection(ctx context.Context) error {
	return mc.cl.Disconnect(ctx)
}

func (md *mongoDB) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)

	return &mongoCollection{collection: collection}
}

func (md *mongoDB) Client() ClientHelper {
	client := md.db.Client()

	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.collection.FindOne(ctx, filter)

	return &mongoSingleResult{sr:singleResult}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (interface{}, error) {
	entryId, err := mc.collection.InsertOne(ctx, document)

	return entryId.InsertedID, err
}

func (mc *mongoCollection) DeleteOne(ctx context.Context, filer interface{}) (int64, error) {
	count, err := mc.collection.DeleteOne(ctx, filer)

	return count.DeletedCount, err
}

func (mc *mongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}) SingleResultHelper {
	singleResult := mc.collection.FindOneAndUpdate(ctx, filter, update, options.FindOneAndUpdate().SetReturnDocument(1))

	return &mongoSingleResult{sr:singleResult}
}

func (mc *mongoCollection) List(ctx context.Context, filter interface{}) (*mongo.Cursor, error) {
	return mc.collection.Find(ctx, filter)
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

