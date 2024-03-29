package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/SowinskiBraeden/DeviceBookingAPI/api"
	"github.com/SowinskiBraeden/DeviceBookingAPI/models"
	"github.com/SowinskiBraeden/DeviceBookingAPI/util"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"

	"github.com/SowinskiBraeden/DeviceBookingAPI/config"
	"github.com/SowinskiBraeden/DeviceBookingAPI/databases"
)

var validate = validator.New()

// App stores the router and db connection so it can be reused
type App struct {
	Router   *mux.Router
	DB       databases.CollectionHelper
	Config   config.Config
	dbHelper databases.DatabaseHelper
}

// New creates a new mux router and all the routes
func (a *App) New() *mux.Router {

	// Detect if system is new and needs default admin
	a.newSystem()

	r := mux.NewRouter()
	cow := Cow{DB: databases.NewCowDatabase(a.dbHelper)}
	device := Device{DB: databases.NewDeviceDatabase(a.dbHelper)}
	// user := User{DB: databases.NewUserDatabase(a.dbHelper)}

	// healthcheck
	r.HandleFunc("/health", healthCheckHandler)

	apiCreate := r.PathPrefix("/api/v1").Subrouter()

	// Data handlers, create, delete, update etc.
	apiCreate.Handle("/cow/{cow_id}", api.Middleware(http.HandlerFunc(cow.CowByObjectIDHandler))).Methods("GET")             // By Object ID not Cow Name
	apiCreate.Handle("/cows", api.Middleware(http.HandlerFunc(cow.CowHandler))).Methods("GET")                               // Returns all cows
	apiCreate.Handle("/cows", api.Middleware(http.HandlerFunc(cow.CowHandlerQuery))).Methods("POST")                         // Returns list of cows based of name query
	apiCreate.Handle("/cows/new", api.Middleware(http.HandlerFunc(cow.NewCowHandler))).Methods("POST")                       // Create new cow
	apiCreate.Handle("/cows/update/{cow_id}", api.Middleware(http.HandlerFunc(cow.UpdateCowHandler))).Methods("POST")        // Update Cow by Object ID
	apiCreate.Handle("/cows/add_device/{cow_id}", api.Middleware(http.HandlerFunc(cow.AddDeviceHandler))).Methods("POST")    // Add Device to cow device list
	apiCreate.Handle("/cows/get_devices/{cow_id}", api.Middleware(http.HandlerFunc(device.GetChildDevices))).Methods("POST") // Returns a list of devices from a given Cow obj
	apiCreate.Handle("cows/bookings/{cow_id}", api.Middleware(http.HandlerFunc(cow.GetBookingsHandler))).Methods("GET")      // Returns all bookings for a given cow

	apiCreate.Handle("/device/{device_id}", api.Middleware(http.HandlerFunc(device.DeviceByObjectIDHandler))).Methods("GET")      // By Object ID not Device Name
	apiCreate.Handle("/devices", api.Middleware(http.HandlerFunc(device.DeviceHandler))).Methods("GET")                           // Returns all devices
	apiCreate.Handle("/devices", api.Middleware(http.HandlerFunc(device.DeviceHandlerQuery))).Methods("POST")                     // Returns list of devices based of name query
	apiCreate.Handle("/devices/new", api.Middleware(http.HandlerFunc(device.NewDeviceHandler))).Methods("POST")                   // create new device
	apiCreate.Handle("/devices/update/{device_id}", api.Middleware(http.HandlerFunc(device.UpdateDeviceHandler))).Methods("POST") // Update Device by Object ID

	// Booking handling
	apiCreate.Handle("/cow/book/{cow_id}", api.Middleware(http.HandlerFunc(cow.BookingHandler))).Methods("POST") // Add booking to cow by ID

	return r
}

func (a *App) Initialize() error {
	client, err := databases.NewClient(&a.Config)
	if err != nil {
		// if we fail to create a new database client, the kill the pod
		zap.S().With(err).Error("failed to create new client")
		return err
	}

	a.dbHelper = databases.NewDatabase(&a.Config, client)
	err = client.Connect()
	if err != nil {
		// if we fail to connect to the database, the kill the pod
		zap.S().With(err).Error("failed to connect to database")
		return err
	}
	zap.S().Info("DeviceBookingAPI has connected to the database")

	// initialize api router
	a.initializeRoutes()
	return nil

}

func (a *App) initializeRoutes() {
	a.Router = a.New()
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, _ := json.Marshal(models.HealthCheckResponse{
		Alive: true,
	})
	_, _ = io.WriteString(w, string(b))
}

func (a *App) newSystem() {
	var DB databases.UserDatabase = databases.NewUserDatabase(a.dbHelper)
	dbResp, err := DB.Find(context.TODO(), bson.M{})
	if err != nil {
		zap.S().With(err).Error("Unable to detect new system: failed to get users")
	}

	if len(dbResp) == 0 {
		fmt.Println("Admin account setup...")

		for {
			defaultAdmin := util.CreateDefaultAdmin(DB)

			if util.Confirm("Are the above credentials correct?") {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()

				_, err := DB.InsertOne(ctx, defaultAdmin)
				if err != nil {
					log.Printf("Failed to create an admin\n")
					break
				}

				log.Printf("Successfully created default admin")
				log.Printf("Your default admin ID is %s", defaultAdmin.Details.UID)
				break
			}
		}
	}
}
