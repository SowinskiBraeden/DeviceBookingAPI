package databases

// go generate: mockery --name CowDatabase

import (
	"context"

	"github.com/SowinskiBraeden/DeviceBookingAPI/models"
)

const cowName = "cows"

// CowDatabase contains the methods to use with the cow database
type CowDatabase interface {
	FindOne(ctx context.Context, filter interface{}) (*models.Cow, error)
	Find(ctx context.Context, filter interface{}) ([]models.Cow, error)
	InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter, document interface{}) (*mongoUpdateResult, error)
}

type cowDatabase struct {
	db DatabaseHelper
}

// NewCowDatabase initialized a new instance of a cow database with the provided db conntection
func NewCowDatabase(db DatabaseHelper) CowDatabase {
	return &cowDatabase{
		db: db,
	}
}

func (c *cowDatabase) FindOne(ctx context.Context, filter interface{}) (*models.Cow, error) {
	cow := &models.Cow{}
	err := c.db.Collection(cowName).FindOne(ctx, filter).Decode(&cow)
	if err != nil {
		return nil, err
	}
	return cow, nil
}

func (c *cowDatabase) Find(ctx context.Context, filter interface{}) ([]models.Cow, error) {
	var cows []models.Cow
	err := c.db.Collection(cowName).Find(ctx, filter).Decode(&cows)
	if err != nil {
		return nil, err
	}
	return cows, nil
}

// Returns the result (document id) and error
func (c *cowDatabase) InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error) {
	result, err := c.db.Collection(cowName).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *cowDatabase) UpdateOne(ctx context.Context, filter, update interface{}) (*mongoUpdateResult, error) {
	result, err := c.db.Collection(cowName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
