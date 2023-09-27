package models

// Device holds the structure for the device collection in mongo
type Device struct {
	ID      string        `json:"_id"    bson:"_id"`    // MongoDB ID
	Details DeviceDetails `json:"device" bson:"Device"` // Details
}

// Device holds the structure for the Device collection in mongo
type DeviceDetails struct {
	Type   string `json:"type"   bson:"Type"`   // eg. Laptop, Ipad, etc
	Name   string `json:"name"   bson:"Name"`   // eg. SULH-LAP-01
	Parent string `json:"parent" bson:"Parent"` // Parent cow ID
}
