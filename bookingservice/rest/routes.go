package rest

import (
	"cloud-microservice-go/lib/msgqueue"
	"cloud-microservice-go/lib/persistence"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func ServeAPI(listenAddr string, database persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter)  {
	r := mux.NewRouter()
	r.Methods("POST").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{eventEmitter, database})

	srv := http.Server{
		Handler: r,
		Addr: listenAddr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout: 1 * time.Second,
	}

	srv.ListenAndServe()
}
