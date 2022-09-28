package databases

// go generate: mockery --name DeviceDatabase

import (
	"context"

	"github.com/SowinskiBraeden/SulliCartShare/models"
)

const deviceName = "devices"

// DeviceDatabase contains the methods to use with the cow database
type DeviceDatabase interface {
	FindOne(ctx context.Context, filter interface{}) (*models.Device, error)
	Find(ctx context.Context, filter interface{}) ([]models.Device, error)
	InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error)
	UpdateOne(ctx context.Context, filter, document interface{}) (*mongoUpdateResult, error)
}

type deviceDatabase struct {
	db DatabaseHelper
}

// NewDeviceDatabase initializes a new instance of a device database with the provided db connection
func NewDeviceDatabase(db DatabaseHelper) DeviceDatabase {
	return &deviceDatabase{
		db: db,
	}
}

func (d *deviceDatabase) FindOne(ctx context.Context, filter interface{}) (*models.Device, error) {
	device := &models.Device{}
	err := d.db.Collection(deviceName).FindOne(ctx, filter).Decode(&device)
	if err != nil {
		return nil, err
	}
	return device, nil
}

func (d *deviceDatabase) Find(ctx context.Context, filter interface{}) ([]models.Device, error) {
	var devices []models.Device
	err := d.db.Collection(deviceName).FindOne(ctx, filter).Decode(&devices)
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (d *deviceDatabase) InsertOne(ctx context.Context, document interface{}) (*mongoInsertOneResult, error) {
	result, err := d.db.Collection(deviceName).InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (d *deviceDatabase) UpdateOne(ctx context.Context, filter, update interface{}) (*mongoUpdateResult, error) {
	result, err := d.db.Collection(deviceName).UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
