package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/SowinskiBraeden/DeviceBookingAPI/config"
	"github.com/SowinskiBraeden/DeviceBookingAPI/databases"
	"github.com/SowinskiBraeden/DeviceBookingAPI/models"
)

type User struct {
	DB databases.UserDatabase
}

// CowByIDHandler returns a cow by ID
func (u User) UserByObjectIDHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_object_id"]

	cID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	dbResp, err := u.DB.FindOne(context.Background(), bson.M{"_id": cID})
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

// TODO: Handle User.Details.UserType properly i.e only SuperUser
// NewUserHandler inserts a new cow into the collection and returns a result and error
func (u User) NewUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var userDetails models.UserDetails // Json data will represent the user details model
	defer cancel()

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&userDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&userDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	newUser := models.User{
		ID:      primitive.NewObjectID().Hex(),
		Details: userDetails,
	}

	result, err := u.DB.InsertOne(ctx, newUser)
	if err != nil {
		config.ErrorStatus("failed to insert user", http.StatusBadRequest, w, err)
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
