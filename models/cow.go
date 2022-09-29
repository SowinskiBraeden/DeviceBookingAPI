package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Cow holds the structure for the cow collection in mongo
type Cow struct {
	ID      string     `json:"_id" bson:"_id"` // MongoDB ID
	Details CowDetails `json:"cow" bson:"Cow"` // Details
}

// BookDetails holds the checkout details
type BookDetails struct {
	Author  string             `json:"author"  bson:"Author"`  // User who booked
	Devices []string           `json:"devices" bson:"Devices"` // Array of device mongo ID's
	Block   string             `json:"block"   bson:"Block"`   // Block that is booked
	Date    primitive.DateTime `json:"date"    bson:"Date"`    // Date this booking occurs
}

// CowDetails holds the structure for the inner cow structure as
// defined in the cow collection in mongo
type CowDetails struct {
	Name        string        `json:"name"        bson:"Name"`        // eg. CA-01
	Collection  string        `json:"collection"  bson:"Collection"`  // eg. Laptop, Ipad, etc
	DeviceTotal int           `json:"deviceTotal" bson:"DeviceTotal"` // # of devices in that cart collection
	Bookings    []BookDetails `json:"bookings"    bson:"Bookings"`    // An array of all active bookings (send top 10)
	Devices     []string      `json:"devices"     bson:"Devices"`     // Array of device mongo ID's
}
