package models

// Device holds the structure for the Device collection in mongo
type Device struct {
	ID     string `json:"ID" bson:"ID"`         // MongoDB ID
	Type   string `json:"Type" bson:"Type"`     // eg. Laptop, Ipdad, etc
	Code   string `json:"Code" bson:"Code"`     // eg. SULH-LAP-01
	Parent string `json:"Parent" bson:"Parent"` // Parent cow Mongo ID
}
