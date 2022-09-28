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
	Author  string             `json:"author" bson:"author"`   // User who booked
	Devices []string           `json:"devices" bson:"devices"` // Array of device mongo ID's
	Block   string             `json:"block" bson:"block"`     // Block that is booked
	Date    primitive.DateTime `json:"date" bson:"date"`       // Date this booking occurs
}

// CowDetails holds the structure for the inner cow structure as
// defined in the cow collection in mongo
type CowDetails struct {
	Name         string        `json:"name" bson:"name"`                 // eg. CA-01
	Collection   string        `json:"collection" bson:"collection"`     // eg. Laptop, Ipad, etc
	TotalDevices int           `json:"totalDevices" bson:"totalDevices"` // # of devices in that cart collection
	Bookings     []BookDetails `json:"bookings" bson:"bookings"`         // An array of all active bookings (send top 10)
	Devices      []string      `json:"devices" bson:"devices"`           // Array of device mongo ID's
}
