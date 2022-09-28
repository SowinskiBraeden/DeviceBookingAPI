package models

// Device holds the structure for the device collection in mongo
type Device struct {
	ID      string        `json:"_id" bson:"_id"`       // MongoDB ID
	Details DeviceDetails `json:"device" bson:"device"` // Details
}

// Device holds the structure for the Device collection in mongo
type DeviceDetails struct {
	Type   string `json:"Type" bson:"Type"`     // eg. Laptop, Ipdad, etc
	Code   string `json:"Code" bson:"Code"`     // eg. SULH-LAP-01
	Parent string `json:"Parent" bson:"Parent"` // Parent cow Mongo ID
}
