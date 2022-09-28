package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cow holds the structure for the cow collection in mongo
type Cow struct {
	ID      string     `json:"_id" bson:"_id"` // MongoDB ID
	Details CowDetails `json:"cow" bson:"cow"` // Details
}

// BookDetails holds the checkout details
type BookDetails struct {
	Author  string             `json:"Author" bson:"Author"`   // User who booked
	Devices []string           `json:"Devices" bson:"Devices"` // Array of device mongo ID's
	Block   string             `json:"Block" bson:"Block"`     // Block that is booked
	Date    primitive.DateTime `json:"Date" bson:"Date"`       // Date this booking occurs
}

// CowDetails holds the structure for the inner cow structure as
// defined in the cow collection in mongo
type CowDetails struct {
	CowCode        string        `json:"CowCode" bson:"CowCode"`               // eg. CA-01
	CollectionType string        `json:"CollectionType" bson:"CollectionType"` // eg. Laptop, Ipad, etc
	TotalDevices   int           `json:"TotalDevices" bson:"TotalDevices"`     // # of devices in that cart collection
	Bookings       []BookDetails `json:"Bookings" bson:"Bookings"`             // An array of all active bookings (send top 10)
	Devices        []string      `json:"Devices" bson:"Devices"`               // Array of device mongo ID's
}
