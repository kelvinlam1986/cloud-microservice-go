package rest

import (
	"cloud-microservice-go/contracts"
	"cloud-microservice-go/lib/msgqueue"
	"cloud-microservice-go/lib/persistence"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type CreateBookingHandler struct {
	eventEmiter msgqueue.EventEmitter
	database persistence.DatabaseHandler
}

type eventRef struct {
	ID string `json:"id"`
	Name string `json:"name,omitempty"`
}

type createBookingRequest struct {
	Seats int `json:"seats"`
	Event eventRef `json:"event"`
}

func (h *CreateBookingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	routeVars := mux.Vars(r)
	eventID, ok := routeVars["eventID"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, "missing route parameter 'eventID")
		return
	}

	eventIDMongo, _ := hex.DecodeString(eventID)
	event, err := h.database.FindEvent(eventIDMongo)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "event %s could not be loaded: %s", eventID, err)
		return
	}

	bookingRequest := createBookingRequest{}
	err = json.NewDecoder(r.Body).Decode(&bookingRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "could not decode JSON body: %s", err)
		return
	}

	if bookingRequest.Seats < 0 {
		w.WriteHeader(400)
		fmt.Fprintf(w, "seat number must be positive (was %d)", bookingRequest.Seats)
		return
	}

	eventIDAsByte, _ := event.ID.MarshalText()
	booking := persistence.Booking{
		Date:time.Now().Unix(),
		EventId: eventIDAsByte,
		Seats: bookingRequest.Seats,
	}

	msg := contracts.EventBookedEvent{
		EventId: event.ID.Hex(),
		UserId: "someUserID",
	}

	h.eventEmiter.Emit(&msg)
	h.database.AddBookingForUser([]byte("someUserID"), booking)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&booking)
}