package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/SowinskiBraeden/SulliCartShare/config"
	"github.com/SowinskiBraeden/SulliCartShare/databases"
	"github.com/SowinskiBraeden/SulliCartShare/models"
)

// Cow exported for testing purposes
type Cow struct {
	DB databases.CowDatabase
}

// CowHandler returns all cows
func (c Cow) CowHandler(w http.ResponseWriter, r *http.Request) {
	dbResp, err := c.DB.Find(context.TODO(), bson.M{})
	if err != nil {
		config.ErrorStatus("failed to get cows", http.StatusNotFound, w, err)
		return
	}

	// If len == 0 then we will just return an empty data object
	if len(dbResp) == 0 {
		dbResp = []models.Cow{}
	}
	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// CowHandlerQuery is the same as CowHanlder, but queries a specific list of objects by Name
func (c Cow) CowHandlerQuery(w http.ResponseWriter, r *http.Request) {
	var query models.Query // Json data will represent the query model

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&query); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	dbResp, err := c.DB.Find(context.TODO(), bson.M{"detials.name": query.Name}) // Search by cow name
	if err != nil {
		config.ErrorStatus("failed to get cows", http.StatusNotFound, w, err)
		return
	}

	// If len == 0 then we will just return an empty data object
	if len(dbResp) == 0 {
		dbResp = []models.Cow{}
	}
	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// CowByIDHandler returns a cow by ID
func (c Cow) CowByIDHandler(w http.ResponseWriter, r *http.Request) {
	cowID := mux.Vars(r)["cow_id"]

	cID, err := primitive.ObjectIDFromHex(cowID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := c.DB.FindOne(context.Background(), bson.M{"_id": cID})
	if err != nil {
		config.ErrorStatus("failed to get cow by ID", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// NewCowHandler inserts a new cow into the collection and returns a result and error
func (c Cow) NewCowHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var cowDetails models.CowDetails // Json data will represent the cow details model
	defer cancel()

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&cowDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&cowDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	newCow := models.Cow{
		ID:      primitive.NewObjectID().Hex(),
		Details: cowDetails,
	}

	result, err := c.DB.InsertOne(ctx, newCow)
	if err != nil {
		config.ErrorStatus("failed to insert cow", http.StatusBadRequest, w, err)
		return
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"result": result}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(b)
}

// UpdateCowHandler gets updates the data for an existing cow and returns a result and error
// This function can only handle updating Name, DeviceTotal, Collection
func (c Cow) UpdateCowHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var newDetails models.CowDetails // Json data will represent the cow details model
	defer cancel()

	cowID := mux.Vars(r)["cow_id"]

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&newDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&newDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	e := reflect.ValueOf(&newDetails).Elem()
	var update bson.M = bson.M{}

	// Only get provided values to update
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varValue := e.Field(i).Interface()
		if varValue != nil && varValue != "" && varName != "Bookings" && varName != "Devices" {
			update["Cow."+varName] = varValue
		}
	}

	dbResp, err := c.DB.UpdateOne(ctx, bson.M{"_id": cowID}, bson.M{"$set": update})
	if err != nil {
		config.ErrorStatus("the cow could not be updated", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func (c Cow) AddDeviceHandler(w http.ResponseWriter, r *http.Request) {
	var newDevice models.NewDeviceToCow

	cowID := mux.Vars(r)["cow_id"]

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&newDevice); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&newDevice); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	update := bson.M{
		"cow.Devices": newDevice,
	}

	dbResp, err := c.DB.UpdateOne(context.TODO(), bson.M{"_id": cowID}, bson.M{"$push": update})
	if err != nil {
		config.ErrorStatus("the device could not be inserted into the cow", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
