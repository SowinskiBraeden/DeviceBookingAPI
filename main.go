package main

import (
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/SowinskiBraeden/DeviceBookingAPI/api/handlers"
	"github.com/SowinskiBraeden/DeviceBookingAPI/config"
)

func main() {
	a := handlers.App{}
	a.Config = *config.New()

	err := a.Initialize() //initialize database and router
	if err != nil {
		zap.S().With(err).Error("error calling initialize")
		return
	}

	zap.S().Infow("DeviceBookingAPI is up and running", "url", a.Config.BaseURL, "port", a.Config.Port)
	log.Fatal(http.ListenAndServe(":"+a.Config.Port, a.Router))
}
