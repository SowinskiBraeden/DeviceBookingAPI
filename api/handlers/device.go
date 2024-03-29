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

	"github.com/SowinskiBraeden/DeviceBookingAPI/config"
	"github.com/SowinskiBraeden/DeviceBookingAPI/databases"
	"github.com/SowinskiBraeden/DeviceBookingAPI/models"
)

type Device struct {
	DB databases.DeviceDatabase
}

// DeviceHandler returns all cows
func (d Device) DeviceHandler(w http.ResponseWriter, r *http.Request) {
	dbResp, err := d.DB.Find(context.TODO(), bson.M{})
	if err != nil {
		config.ErrorStatus("failed to get devices", http.StatusNotFound, w, err)
		return
	}

	// If len == 0 then we will just return an empty data object
	if len(dbResp) == 0 {
		dbResp = []models.Device{}
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// DeviceHandlerQuery is the same as DeviceHandler, but queries a specific list of objects by Name
func (d Device) DeviceHandlerQuery(w http.ResponseWriter, r *http.Request) {
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

	dbResp, err := d.DB.Find(context.TODO(), bson.M{"detials.name": query.Name}) // Search by device name
	if err != nil {
		config.ErrorStatus("failed to get cow(s)", http.StatusNotFound, w, err)
		return
	}

	// If len == 0 then we will just return an empty data object
	if len(dbResp) == 0 {
		dbResp = []models.Device{}
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": dbResp}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

// DeviceByIDHandler returns a cow by ID
func (d Device) DeviceByObjectIDHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := mux.Vars(r)["device_id"]

	dbResp, err := d.DB.FindOne(context.Background(), bson.M{"_id": deviceID})
	if err != nil {
		config.ErrorStatus("failed to get device by ObjectID", http.StatusNotFound, w, err)
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

// NewDeviceHandler inserts a new cow into the collection and returns a result and error
func (d Device) NewDeviceHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var deviceDetails models.DeviceDetails // Json data will represent the device details model
	defer cancel()

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&deviceDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&deviceDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	newDevice := models.Device{
		ID:      primitive.NewObjectID().Hex(),
		Details: deviceDetails,
	}

	result, err := d.DB.InsertOne(ctx, newDevice)
	if err != nil {
		config.ErrorStatus("failed to insert device", http.StatusInternalServerError, w, err)
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

// UpdateDeviceHandler gets updates the data for an existing cow and returns a result and error
func (d Device) UpdateDeviceHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var newDetails models.DeviceDetails // Json data will represent the cow details model
	defer cancel()

	deviceID := mux.Vars(r)["device_id"]

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
		if varValue != nil && varValue != "" {
			update["Device."+varName] = varValue
		}
	}

	dbResp, err := d.DB.UpdateOne(ctx, bson.M{"_id": deviceID}, bson.M{"$set": update})
	if err != nil {
		config.ErrorStatus("the device could not be updated", http.StatusNotFound, w, err)
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

// Is able to get child devices with a prodived Cow Model
func (d Device) GetChildDevices(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	cowID := mux.Vars(r)["cow_id"]
	defer cancel()

	devices, _ := d.DB.Find(ctx, bson.M{"Device.Parent": cowID})

	// If there is no devices from the query, return empty device array.
	if len(devices) == 0 {
		devices = []models.Device{}
	}

	b, err := json.Marshal(models.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"result": devices}})
	if err != nil {
		config.ErrorStatus("failed to marshal response", http.StatusInternalServerError, w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}
