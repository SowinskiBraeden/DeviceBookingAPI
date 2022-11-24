package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/SowinskiBraeden/SulliCartShare/config"
	"github.com/SowinskiBraeden/SulliCartShare/models"
	"github.com/gorilla/mux"
	"github.com/thanhpk/randstr"
)

func (c Cow) BookingHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var bookingDetails models.BookDetails // Json data will represent the cow details model
	defer cancel()

	// validate the request body
	if err := json.NewDecoder(r.Body).Decode(&bookingDetails); err != nil {
		config.ErrorStatus("failed to unpack request body", http.StatusInternalServerError, w, err)
		return
	}

	// use the validator library to validate required fields
	if validationErr := validate.Struct(&bookingDetails); validationErr != nil {
		config.ErrorStatus("invalid request body", http.StatusBadRequest, w, validationErr)
		return
	}

	cowID := mux.Vars(r)["cow_id"]

	cID, err := primitive.ObjectIDFromHex(cowID)
	if err != nil {
		config.ErrorStatus("failed to get objectID from Hex", http.StatusBadRequest, w, err)
		return
	}

	bookingDetails.ID = fmt.Sprintf("%s.%s", cowID, randstr.Hex(16))

	dbResp, err := c.DB.UpdateOne(ctx, bson.M{"_id": cID}, bson.M{"$push": bson.M{"Cow.Bookings": bookingDetails}})
	if err != nil {
		config.ErrorStatus("the booking could not be added to the cow", http.StatusNotFound, w, err)
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
