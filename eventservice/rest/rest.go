package rest

import (
	"cloud-microservice-go/lib/msgqueue"
	"cloud-microservice-go/lib/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

func ServeAPI(endpoint string, databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) error {
	handler := newEventHandler(databaseHandler, eventEmitter)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/events").Subrouter()
	eventsRouter.Methods("GET").PathPrefix("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsRouter.Methods("GET").PathPrefix("").HandlerFunc(handler.AllEventHandler)
	eventsRouter.Methods("GET").PathPrefix("/{eventID}").HandlerFunc(handler.OneEventHandler)
	eventsRouter.Methods("POST").PathPrefix("").HandlerFunc(handler.NewEventHandler)

	locationRouter := r.PathPrefix("/locations").Subrouter()
	locationRouter.Methods("POST").PathPrefix("").HandlerFunc(handler.NewLocationHandler)
	locationRouter.Methods("GET").PathPrefix("").HandlerFunc(handler.AllLocationsHandler)
	
	return http.ListenAndServe(endpoint, r)
}
