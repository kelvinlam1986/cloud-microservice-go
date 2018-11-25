package rest

import (
	"cloud-microservice-go/lib/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

func ServeAPI(endpoint string, databaseHandler persistence.DatabaseHandler) error {
	handler := New(databaseHandler)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("events").Subrouter()
	eventsRouter.Methods("GET").PathPrefix("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsRouter.Methods("GET").PathPrefix("").HandlerFunc(handler.AllEventHandler)
	eventsRouter.Methods("POST").PathPrefix("").HandlerFunc(handler.NewEventHandler)
	return http.ListenAndServe(endpoint, r)
}
