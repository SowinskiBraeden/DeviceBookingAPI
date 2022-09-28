package databases

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/SowinskiBraeden/SulliCartShare/config"
)

// DatabaseHelper contains the collection and client to be used to access the methods
// defined below
type DatabaseHelper interface {
	Collection(name string) CollectionHelper
	Client() ClientHelper
}

// CollectionHelper contains all the methods defined for collection in this project
type CollectionHelper interface {
	FindOne(context.Context, interface{}) SingleResultHelper
	Find(context.Context, interface{}) CursorHelper
	InsertOne(context.Context, interface{}) (mongoInsertOneResult, error)
	UpdateOne(context.Context, interface{}, interface{}) (mongoUpdateResult, error)
}

// SingleResultHelper contains a single method to decode the result
type SingleResultHelper interface {
	Decode(v interface{}) error
}

// CursorHelper contains a method to decode the cursor
type CursorHelper interface {
	Decode(v interface{}) error
}

// ClientHelper defined to help at client creation inside main.go
type ClientHelper interface {
	Database(string) DatabaseHelper
	Connect() error
	StartSession() (mongo.Session, error)
}

type mongoClient struct {
	cl *mongo.Client
}

type mongoDatabase struct {
	db *mongo.Database
}

type mongoCollection struct {
	coll *mongo.Collection
}

type mongoSingleResult struct {
	sr *mongo.SingleResult
}

type mongoCursor struct {
	cr *mongo.Cursor
}

type mongoInsertOneResult struct {
	ir *mongo.InsertOneResult
}

type mongoUpdateResult struct {
	Ur *mongo.UpdateResult // capital to export field to ignore error in ../api/handlers/cow.go Ln 133 Col 26
}

type mongoSession struct {
	mongo.Session
}

// NewClient uses the values from the config and returns a mongo client
func NewClient(conf *config.Config) (ClientHelper, error) {
	c, err := mongo.NewClient(options.Client().ApplyURI(conf.URL))

	return &mongoClient{cl: c}, err
}

// NewDatabase uses the client from NewClient and sets the database name
func NewDatabase(conf *config.Config, client ClientHelper) DatabaseHelper {
	return client.Database(conf.DatabaseName)
}

func (mc *mongoClient) Database(dbName string) DatabaseHelper {
	db := mc.cl.Database(dbName)
	return &mongoDatabase{db: db}
}

func (mc *mongoClient) StartSession() (mongo.Session, error) {
	session, err := mc.cl.StartSession()
	return &mongoSession{session}, err
}

func (mc *mongoClient) Connect() error {
	return mc.cl.Connect(context.TODO()) // use context.TODO() instead of nil cause good practice ¯\_(ツ)_/¯
}

func (md *mongoDatabase) Collection(colName string) CollectionHelper {
	collection := md.db.Collection(colName)
	return &mongoCollection{coll: collection}
}

func (md *mongoDatabase) Client() ClientHelper {
	client := md.db.Client()
	return &mongoClient{cl: client}
}

func (mc *mongoCollection) FindOne(ctx context.Context, filter interface{}) SingleResultHelper {
	singleResult := mc.coll.FindOne(ctx, filter)
	return &mongoSingleResult{sr: singleResult}
}

func (mc *mongoCollection) Find(ctx context.Context, filter interface{}) CursorHelper {
	cursor, _ := mc.coll.Find(ctx, filter)
	return &mongoCursor{cr: cursor}
}

func (mc *mongoCollection) InsertOne(ctx context.Context, document interface{}) (mongoInsertOneResult, error) {
	insertOneResult, err := mc.coll.InsertOne(ctx, document)
	if err != nil {
		return mongoInsertOneResult{}, err
	}
	return mongoInsertOneResult{ir: insertOneResult}, nil
}

func (mc *mongoCollection) UpdateOne(ctx context.Context, filter, update interface{}) (mongoUpdateResult, error) {
	updateOneResult, err := mc.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return mongoUpdateResult{}, err
	}
	return mongoUpdateResult{Ur: updateOneResult}, nil
}

func (sr *mongoSingleResult) Decode(v interface{}) error {
	return sr.sr.Decode(v)
}

func (cr *mongoCursor) Decode(v interface{}) error {
	return cr.All(context.Background(), v)
}

func (cr *mongoCursor) All(ctx context.Context, results interface{}) error {
	return cr.cr.All(ctx, results)
}
