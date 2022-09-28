package models

// Device holds the structure for the device collection in mongo
type Device struct {
	ID      string        `json:"_id" bson:"_id"`         // MongoDB ID
	Details DeviceDetails `json:"details" bson:"details"` // Details
}

// Device holds the structure for the Device collection in mongo
type DeviceDetails struct {
	Type   string `json:"type" bson:"type"`     // eg. Laptop, Ipdad, etc
	Name   string `json:"name" bson:"name"`     // eg. SULH-LAP-01
	Parent string `json:"parent" bson:"parent"` // Parent cow Mongo ID
}
