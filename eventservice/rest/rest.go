package rest

import (
	"cloud-microservice-go/lib/msgqueue"
	"cloud-microservice-go/lib/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

func ServeAPI(endpoint, tlsendpoint string, databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := newEventHandler(databaseHandler, eventEmitter)
	r := mux.NewRouter()
	eventsRouter := r.PathPrefix("/events").Subrouter()
	eventsRouter.Methods("GET").PathPrefix("/{SearchCriteria}/{search}").HandlerFunc(handler.FindEventHandler)
	eventsRouter.Methods("GET").PathPrefix("").HandlerFunc(handler.AllEventHandler)
	eventsRouter.Methods("POST").PathPrefix("").HandlerFunc(handler.NewEventHandler)

	httpErrorChan := make(chan error)
	httptlsErrorChan := make(chan error)

	go func() { httptlsErrorChan <- http.ListenAndServeTLS(tlsendpoint, "cert.pem", "key.pem", r) }()
	go func() { httpErrorChan <- http.ListenAndServe(endpoint, r) }()

	return httpErrorChan, httptlsErrorChan
}
