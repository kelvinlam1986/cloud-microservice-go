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
	"strings"
	"time"
)

type eventServiceHandler struct {
	dbHandler persistence.DatabaseHandler
	eventEmiiter msgqueue.EventEmitter
}

func newEventHandler(databaseHandler persistence.DatabaseHandler, eventEmiter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbHandler: databaseHandler,
		eventEmiiter: eventEmiter,
	}
}

func (eh *eventServiceHandler) FindEventHandler(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		fmt.Fprint(w, `No search criteria found, you can either search by id via /id/4
						to search by name via /name/coldplayconcert`)
		return
	}

	searchKey, ok := vars["search"]
	if !ok {
		fmt.Fprint(w, `No search keys found, you can either search by id via /id/4
						to search by name via /name/coldplayconcert`)
		return
	}

	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
	case "name":
		event, err = eh.dbHandler.FindEventByName(searchKey)
	case "id":
		id, err := hex.DecodeString(searchKey)
		if err == nil {
			event, err = eh.dbHandler.FindEvent(id)
		}
	}

	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "Error occured %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) AllEventHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("All events we are here")
	events, err := eh.dbHandler.FindAllAvailableEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying to find all available events %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error occured while trying encode events to JSON %s", err)
	}
}

func (eh *eventServiceHandler) NewEventHandler(w http.ResponseWriter, r *http.Request)  {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while decoding event data %s", err)
		return
	}

	id, err := eh.dbHandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "error occured while persisting event %s", err)
		return
	}

	msg := contracts.EventCreatedEvent{
		ID: hex.EncodeToString(id),
		Name: event.Name,
		LocationID: string(event.Location.ID),
		Start: time.Unix(event.StartDate,0),
		End: time.Unix(event.EndDate, 0),
	}

	eh.eventEmiiter.Emit(&msg)
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&event)

	fmt.Fprint(w, id)
}