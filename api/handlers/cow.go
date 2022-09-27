package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-playground/validator/v10"

	"github.com/SowinskiBraeden/SulliCartShare/config"
	"github.com/SowinskiBraeden/SulliCartShare/databases"
	"github.com/SowinskiBraeden/SulliCartShare/models"
)

var validate = validator.New()

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
	b, err := json.Marshal(dbResp)
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

	b, err := json.Marshal(dbResp)
	if err != nil {
		config.ErrorStatus("failed to marshal responce", http.StatusInternalServerError, w, err)
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

	//validate the request body
	if err := json.NewDecoder(r.Body).Decode(&cowDetails); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"error": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	//use the validator library to validate required fields
	if validationErr := validate.Struct(&cowDetails); validationErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := models.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"error": validationErr.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	fmt.Print(cowDetails)

	newCow := models.Cow{
		ID:      primitive.NewObjectID().Hex(),
		Details: cowDetails,
	}

	result, err := c.DB.InsertOne(ctx, newCow)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response := models.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"error": err.Error()}}
		json.NewEncoder(w).Encode(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := models.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"result": result}}
	json.NewEncoder(w).Encode(response)
}

// UpdateCowHandler gets updates the data for an existing cow and returns a result and error
func (c Cow) UpdateCowHandler(w http.ResponseWriter, r *http.Request) {
	cowID := mux.Vars(r)["cow_id"]

	// TODO: Collect data to update from passed json
	update := bson.M{
		"$set": bson.M{
			"cow.CowCode": "123",
		},
	}

	dbResp, err := c.DB.UpdateOne(context.Background(), bson.M{"_id": cowID}, update)
	if err != nil {
		config.ErrorStatus("failed to update cow by ID", http.StatusNotFound, w, err)
		return
	}

	b, err := json.Marshal(dbResp) // TODO: get rid of this warning cause its bothering me
	if err != nil {
		config.ErrorStatus("failed ot marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
