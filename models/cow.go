package models

// Cow holds the structure for the cow collection in mongo
type Cow struct {
	ID      string     `json:"_id" bson:"_id"`
	Details CowDetails `json:"cow" bson:"cow"`
}

// CowDetails holds teh structure for the inner cow structure as
// defined in the cow collection in mongo
type CowDetails struct {
	CowCode string `json:"CowCode" bson:"CowCode"`
}
